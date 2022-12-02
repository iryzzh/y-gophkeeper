//nolint:cyclop
package sqlite

import (
	"context"
	"database/sql"

	"github.com/iryzzh/y-gophkeeper/internal/models"
	"github.com/iryzzh/y-gophkeeper/internal/store"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type ItemRepository struct {
	db *sql.DB
}

func (r *ItemRepository) Create(ctx context.Context, item *models.Item) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if item.Meta == "" {
		return store.ErrItemMetaIsRequired
	}

	if item.ItemData != nil && item.ItemData.Data != nil {
		if err = tx.QueryRowContext(ctx,
			`insert into items_data (data) values (?) returning id`,
			item.ItemData.Data).Scan(&item.DataID); err != nil || item.DataID == 0 {
			vErr, ok := errors.Cause(err).(sqlite3.Error)
			if ok && vErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return store.ErrItemExists
			}

			return errors.Wrap(err, store.ErrItemDataCreateFailed.Error())
		}

		item.ItemData.ID = item.DataID
	}

	if err = tx.QueryRowContext(ctx,
		`insert into items (user_id, meta, data_id, data_type) values (?, ?, ?, ?) returning id, created_at`,
		item.UserID, item.Meta, item.DataID, item.DataType).Scan(&item.ID, &item.CreatedAt); err != nil {
		vErr, ok := errors.Cause(err).(sqlite3.Error)
		if ok && vErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return store.ErrItemExists
		}

		return errors.Wrap(err, store.ErrItemCreateFailed.Error())
	}

	return tx.Commit()
}

func (r *ItemRepository) FindByID(ctx context.Context, userID string, id int) (*models.Item, error) {
	item := &models.Item{}
	itemData := &models.ItemData{}
	err := r.db.QueryRowContext(ctx,
		`select items.id, items.user_id, items.meta, items.data_id, items.data_type, items.created_at,
	       			items.updated_at, idt.id, idt.data
					from items
					left join items_data idt on idt.id = items.data_id				
					where items.user_id = $1 and items.id = $2				
					order by items.id`,
		userID, id).
		Scan(&item.ID, &item.UserID, &item.Meta, &item.DataID, &item.DataType, &item.CreatedAt, &item.UpdatedAt, &itemData.ID,
			&itemData.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrItemNotFound
		}
		return nil, err
	}

	item.ItemData = itemData

	return item, nil
}

func (r *ItemRepository) FindByMetaName(ctx context.Context, userID string, metaName string) (*models.Item, error) {
	item := &models.Item{}
	itemData := &models.ItemData{}
	err := r.db.QueryRowContext(ctx,
		`select items.id, items.user_id, items.meta, items.data_id, items.data_type, items.created_at,
	       			items.updated_at, idt.id, idt.data
					from items
					left join items_data idt on idt.id = items.data_id				
					where items.user_id = $1 and items.meta = $2				
					order by items.id`,
		userID, metaName).
		Scan(&item.ID, &item.UserID, &item.Meta, &item.DataID, &item.DataType, &item.CreatedAt,
			&item.UpdatedAt, &itemData.ID, &itemData.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrItemNotFound
		}
		return nil, err
	}

	item.ItemData = itemData

	return item, nil
}

func (r *ItemRepository) FindByUserID(ctx context.Context, userID string, limit, offset int) (*models.Items, error) {
	rows, err := r.db.QueryContext(ctx,
		`select items.id, items.user_id, items.meta, items.data_id, items.data_type,  items.created_at,
       			items.updated_at, idt.id, idt.data,
       			(select count(*) from items where user_id = $1)
				from items
				left join items_data idt on idt.id = items.data_id				
				where user_id = $1
				group by items.id
				order by items.id
				limit $2 offset $3`,
		userID,
		limit,
		offset)
	if err != nil {
		return nil, errors.Wrap(err, store.ErrItemNotFound.Error())
	}
	defer func() { _ = rows.Close() }()

	var total int
	var items []*models.Item
	for rows.Next() {
		item := &models.Item{}
		itemData := &models.ItemData{}
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Meta,
			&item.DataID,
			&item.DataType,
			&item.CreatedAt,
			&item.UpdatedAt,
			&itemData.ID,
			&itemData.Data,
			&total,
		)
		if err != nil {
			panic("item FindByUserID sql scan error: " + err.Error())
		}
		item.ItemData = itemData
		items = append(items, item)
	}

	if len(items) > 0 {
		return &models.Items{
			Meta: models.Meta{TotalItems: total},
			Data: items,
		}, rows.Err()
	}

	return nil, store.ErrItemNotFound
}

func (r *ItemRepository) Update(ctx context.Context, item *models.Item) error {
	if item.Meta == "" {
		return store.ErrItemMetaIsRequired
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if item.ItemData != nil {
		if item.ItemData.ID == 0 {
			return store.ErrItemDataInvalidID
		}
		res, _ := tx.ExecContext(ctx,
			`update items_data set data = $1 
                  	where id = (select data_id from items where user_id = $2 and data_id = $3)`,
			item.ItemData.Data, item.UserID, item.DataID)
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return store.ErrItemDataNotFound
		}
	} else {
		item.ItemData = &models.ItemData{}
	}

	if item.ID == 0 {
		return store.ErrItemInvalidID
	}
	res, err := tx.ExecContext(ctx,
		`update items set meta = $1, data_id = $2, data_type = $3 where id = $4 and user_id = $5`,
		item.Meta, item.ItemData.ID, item.DataType, item.ID, item.UserID)
	if err != nil {
		return errors.Wrap(err, store.ErrItemUpdateFailed.Error())
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return store.ErrItemNotFound
	}

	return tx.Commit()
}

func (r *ItemRepository) Delete(ctx context.Context, item *models.Item) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx,
		`delete from items_data where id = (select data_id from items where user_id = $1 and id = $2)`,
		item.UserID, item.ID)
	if err != nil {
		return errors.Wrap(err, store.ErrItemDataDeleteFailed.Error())
	}

	_, err = tx.ExecContext(ctx,
		`delete from items where user_id = $1 and id = $2`,
		item.UserID, item.ID)
	if err != nil {
		return errors.Wrap(err, store.ErrItemDeleteFailed.Error())
	}

	return tx.Commit()
}
