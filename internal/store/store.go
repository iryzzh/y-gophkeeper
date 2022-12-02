package store

import (
	"context"

	"github.com/iryzzh/y-gophkeeper/internal/models"
)

// Store is an interface that presents interfaces to users and data.
type Store interface {
	Status
	User() UserRepository
	Item() ItemRepository
	Close() error
}

// Status represents ways to check the status of the database.
type Status interface {
	Ping() error
	IsUsersExist() (bool, error)
	IsItemsExist() (bool, error)
}

// UserRepository represents ways to interact with users in the database.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByLogin(ctx context.Context, login string) (*models.User, error)
	FindByID(ctx context.Context, userID string) (*models.User, error)
}

// ItemRepository represents ways to interact with items in the database.
type ItemRepository interface {
	Create(ctx context.Context, item *models.Item) error
	FindByID(ctx context.Context, userID string, id int) (*models.Item, error)
	FindByMetaName(ctx context.Context, userID string, metaName string) (*models.Item, error)
	FindByUserID(ctx context.Context, userID string, limit, offset int) (*models.Items, error)
	Update(ctx context.Context, item *models.Item) error
	Delete(ctx context.Context, item *models.Item) error
}
