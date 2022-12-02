package item

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/iryzzh/y-gophkeeper/internal/models"
	"github.com/iryzzh/y-gophkeeper/internal/store"
)

var (
	// ErrItemNotFound is returned when the item is not found.
	ErrItemNotFound = errors.New("item not found")
	// ErrIncorrectItemID is returned when the received item id is incorrect.
	ErrIncorrectItemID = errors.New("incorrect item id")
)

// Service is the service responsible for processing items.
type Service struct {
	store store.Store
}

// NewService creates a new item service.
func NewService(s store.Store) *Service {
	return &Service{store: s}
}

// Create creates a new item in the database.
func (s *Service) Create(ctx context.Context, item *models.Item) error {
	err := s.store.Item().Create(ctx, item)

	return err
}

// FindByID returns the item with the given id.
func (s *Service) FindByID(ctx context.Context, userID, id string) (*models.Item, error) {
	var i int
	if _, err := fmt.Sscan(id, &i); err != nil {
		return nil, ErrIncorrectItemID
	}

	item, err := s.store.Item().FindByID(ctx, userID, i)
	if errors.Is(err, store.ErrItemNotFound) {
		return nil, ErrItemNotFound
	}

	return item, err
}

// FindByMetaName returns the item with the given meta name.
func (s *Service) FindByMetaName(ctx context.Context, userID, metaName string) (*models.Item, error) {
	item, err := s.store.Item().FindByMetaName(ctx, userID, metaName)
	if errors.Is(err, store.ErrItemNotFound) {
		return nil, ErrItemNotFound
	}

	return item, err
}

// FindByUserID returns a set of items.
func (s *Service) FindByUserID(ctx context.Context, userID, limit, offset string) (*models.Items, error) {
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	if intLimit == 0 {
		intLimit = 1000
	}

	if intOffset == 0 {
		intOffset = 0
	}

	items, err := s.store.Item().FindByUserID(ctx, userID, intLimit, intOffset)
	if errors.Is(err, store.ErrItemNotFound) {
		return nil, ErrItemNotFound
	}

	return items, err
}

// Update updates the item in the database.
func (s *Service) Update(ctx context.Context, item *models.Item) error {
	if item.ItemData == nil {
		item.ItemData = nil
	}
	err := s.store.Item().Update(ctx, item)
	return err
}

// Delete deletes the item from the database.
func (s *Service) Delete(ctx context.Context, item *models.Item) error {
	err := s.store.Item().Delete(ctx, item)

	return err
}
