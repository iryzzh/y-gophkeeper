package sqlite

import (
	"context"
	"database/sql"
	"testing"

	"github.com/mattn/go-sqlite3"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
	"github.com/iryzzh/gophkeeper/internal/models"
	"github.com/iryzzh/gophkeeper/internal/rand"
	"github.com/iryzzh/gophkeeper/internal/store"
	"github.com/stretchr/testify/require"
)

func TestItemRepository_Create(t *testing.T) {
	tests := []struct {
		name     string
		wantErr  error
		items    []*models.Item
		extended sqlite3.ErrNoExtended
	}{
		{
			name:    "should create a new item",
			wantErr: nil,
			items: []*models.Item{
				{
					UserID: "e030240c-2906-467d-a280-16f9cd197022",
					Meta:   "random/meta",
					ItemData: &models.ItemData{
						Data: []byte("abcd/efg"),
					},
				},
			},
		},
		{
			name:     "item creation should fail - empty meta",
			wantErr:  store.ErrItemCreateFailed,
			extended: sqlite3.ErrConstraintCheck,
			items: []*models.Item{
				{
					UserID: "e030240c-2906-467d-a280-16f9cd197022",
					Meta:   "",
					ItemData: &models.ItemData{
						Data: []byte("abcd/efg1"),
					},
				},
			},
		},
		{
			name:     "item creation should fail - unique",
			wantErr:  store.ErrItemCreateFailed,
			extended: sqlite3.ErrConstraintUnique,
			items: []*models.Item{
				{
					UserID: "e030240c-2906-467d-a280-16f9cd197022",
					Meta:   "random/meta",
					ItemData: &models.ItemData{
						Data: []byte("abcd/efg1"),
					},
				},
				{
					UserID: "e030240c-2906-467d-a280-16f9cd197022",
					Meta:   "random/meta",
					ItemData: &models.ItemData{
						Data: []byte("abcd/efg1"),
					},
				},
			},
		},
		{
			name:     "item data creation should fail - empty data",
			wantErr:  store.ErrItemDataCreateFailed,
			extended: sqlite3.ErrConstraintNotNull,
			items: []*models.Item{
				{
					UserID: "e030240c-2906-467d-a280-16f9cd197022",
					Meta:   "random/meta/123",
					ItemData: &models.ItemData{
						Data: nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ItemRepository{
				db: setupStore(t),
			}
			defer func() { _ = r.db.Close() }()
			var err error
			for _, item := range tt.items {
				err = r.Create(context.Background(), item)
			}
			require.Conditionf(t, func() (success bool) {
				if tt.wantErr == nil && err == nil {
					return true
				}

				if assert.ErrorContains(t, err, tt.wantErr.Error()) {
					switch vErr := errors.Cause(err).(type) {
					case sqlite3.Error:
						return vErr.ExtendedCode == tt.extended
					default:
						t.Fatalf("unimplemented error: %v", vErr)
					}
				}

				return false
			}, "Create failed: expected %v, got %v", tt.wantErr, err)
		})
	}
}

func generateTestItems(t *testing.T, db *sql.DB, userID string, count int) []*models.Item {
	r := &ItemRepository{
		db: db,
	}

	var items []*models.Item
	for i := 0; i < count; i++ {
		item := sampleItem(t, userID)

		if err := r.Create(context.Background(), item); err != nil {
			t.Fatal(err)
		}
		items = append(items, item)
	}

	return items
}

func TestItemRepository_FindByUserID(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		wantErr error
		count   int
	}{
		{
			name:    "ok",
			userID:  uuid.NewString(),
			count:   111,
			wantErr: nil,
		},
		{
			name:    "not found",
			count:   0,
			wantErr: store.ErrItemNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ItemRepository{
				db: setupStore(t),
			}
			defer func() { _ = r.db.Close() }()

			genItems := generateTestItems(t, r.db, tt.userID, tt.count)
			items, pages, err := r.FindByUserID(context.Background(), tt.userID, 1)
			require.Equal(t, tt.wantErr == nil, err == nil)
			if tt.wantErr == nil {
				var allItems []*models.Item
				allItems = append(allItems, items...)
				for page := 2; page <= pages; page++ {
					items, _, err = r.FindByUserID(context.Background(), tt.userID, page)
					require.NoError(t, err)
					allItems = append(allItems, items...)
				}
				require.Equal(t, tt.count, len(allItems))
				require.Equal(t, allItems, genItems)
			}
		})
	}
}

func TestItemRepository_Update(t *testing.T) {
	db := setupStore(t)
	defer func() { _ = db.Close() }()
	tests := []struct {
		name     string
		wantErr  error
		item     *models.Item
		extended sqlite3.ErrNoExtended
	}{
		{
			name:    "ok",
			wantErr: nil,
			item: func() *models.Item {
				item := sampleItem(t, uuid.NewString())
				r := &ItemRepository{
					db: db,
				}
				if err := r.Create(context.Background(), item); err != nil {
					t.Fatal(err)
				}
				return item
			}(),
		},
		{
			name:    "item not found",
			wantErr: store.ErrItemNotFound,
			item: func() *models.Item {
				item := sampleItem(t, uuid.NewString())
				r := &ItemRepository{
					db: db,
				}
				if err := r.Create(context.Background(), item); err != nil {
					t.Fatal(err)
				}
				item.ID = 1234
				return item
			}(),
		},
		{
			name:    "data not found",
			wantErr: store.ErrItemDataNotFound,
			item: func() *models.Item {
				item := sampleItem(t, uuid.NewString())
				r := &ItemRepository{
					db: db,
				}
				if err := r.Create(context.Background(), item); err != nil {
					t.Fatal(err)
				}
				item.DataID = 1234

				return item
			}(),
		},
		{
			name:    "item data invalid id",
			wantErr: store.ErrItemDataInvalidID,
			item: func() *models.Item {
				item := sampleItem(t, uuid.NewString())
				r := &ItemRepository{
					db: db,
				}
				if err := r.Create(context.Background(), item); err != nil {
					t.Fatal(err)
				}
				item.ItemData.ID = 0
				return item
			}(),
		},
		{
			name:    "item invalid id",
			wantErr: store.ErrItemInvalidID,
			item: func() *models.Item {
				item := sampleItem(t, uuid.NewString())
				r := &ItemRepository{
					db: db,
				}
				if err := r.Create(context.Background(), item); err != nil {
					t.Fatal(err)
				}
				item.ID = 0
				return item
			}(),
		},
		{
			name:     "item empty meta",
			wantErr:  store.ErrItemUpdateFailed,
			extended: sqlite3.ErrConstraintCheck,
			item: func() *models.Item {
				item := sampleItem(t, uuid.NewString())
				r := &ItemRepository{
					db: db,
				}
				if err := r.Create(context.Background(), item); err != nil {
					t.Fatal(err)
				}
				item.Meta = ""
				return item
			}(),
		},
		{
			name: "item empty data",
			item: func() *models.Item {
				item := sampleItem(t, uuid.NewString())
				r := &ItemRepository{
					db: db,
				}
				if err := r.Create(context.Background(), item); err != nil {
					t.Fatal(err)
				}
				item.ItemData = nil
				return item
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ItemRepository{
				db: db,
			}
			err := r.Update(context.Background(), tt.item)
			require.Conditionf(t, func() (success bool) {
				if tt.wantErr == nil && err == nil {
					return true
				}

				if assert.ErrorContains(t, err, tt.wantErr.Error()) {
					switch vErr := errors.Cause(err).(type) {
					case sqlite3.Error:
						return vErr.ExtendedCode == tt.extended
					default:
						return vErr.Error() == tt.wantErr.Error()
					}
				}

				return false
			}, "Update failed: expected %v, got %v", tt.wantErr, err)
		})
	}
}

func sampleItem(t *testing.T, userID string) *models.Item {
	t.Helper()
	return &models.Item{
		UserID: userID,
		Meta:   rand.String(10),
		ItemData: &models.ItemData{
			Data: []byte(rand.String(10)),
		},
	}
}

func TestItemRepository_Delete(t *testing.T) {
	db := setupStore(t)
	tests := []struct {
		name    string
		wantErr error
		item    *models.Item
	}{
		{
			name:    "ok",
			wantErr: nil,
			item: func() *models.Item {
				item := sampleItem(t, uuid.NewString())
				r := &ItemRepository{
					db: db,
				}
				if err := r.Create(context.Background(), item); err != nil {
					t.Fatal(err)
				}
				return item
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ItemRepository{
				db: db,
			}
			defer func() { _ = r.db.Close() }()
			err := r.Delete(context.Background(), tt.item)
			require.Equal(t, tt.wantErr == nil, err == nil)
		})
	}
}
