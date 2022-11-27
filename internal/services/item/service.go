package item

import (
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"

	"github.com/iryzzh/gophkeeper/internal/models"
	"github.com/iryzzh/gophkeeper/internal/store"
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
	log.Printf("item create svc err: %v", err)

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

// FindByUserID returns a set of items.
func (s *Service) FindByUserID(ctx context.Context, userID string, page any) ([]*models.Item, int, error) {
	// already checked in `pagination` func
	p, _ := page.(int)

	items, pages, err := s.store.Item().FindByUserID(ctx, userID, p)
	if errors.Is(err, store.ErrItemNotFound) {
		return nil, 0, ErrItemNotFound
	}

	return items, pages, err
}

// Update updates the item in the database.
func (s *Service) Update(ctx context.Context, item *models.Item) error {
	if item.ItemData == nil {
		item.ItemData = nil
	}
	err := s.store.Item().Update(ctx, item)
	log.Printf("item update svc err: %v", err)
	return err
}

// Delete deletes the item from the database.
func (s *Service) Delete(ctx context.Context, item *models.Item) error {
	err := s.store.Item().Delete(ctx, item)
	log.Printf("item delete svc err: %v", err)

	return err
}
