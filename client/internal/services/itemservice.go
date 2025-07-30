package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rycln/gokeep/shared/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// Interfaces for storage operations (segregated by function)
type itemStorer interface {
	Add(context.Context, *models.ItemInfo, []byte) error
}

type itemGetter interface {
	ListByUser(context.Context, models.UserID) ([]models.ItemInfo, error)
	GetContent(context.Context, models.ItemID) ([]byte, error)
}

type itemDeleter interface {
	DeleteItem(context.Context, models.ItemID) error
}

type itemUpdater interface {
	UpdateItem(context.Context, *models.ItemInfo, []byte) error
}

// itemStorage combines all storage operations interfaces
type itemStorage interface {
	itemStorer
	itemGetter
	itemDeleter
	itemUpdater
}

// ItemService handles business logic for item operations
type ItemService struct {
	storage itemStorage
}

// NewItemService creates a new ItemService instance
func NewItemService(storage itemStorage) *ItemService {
	return &ItemService{
		storage: storage,
	}
}

// Add creates a new item with generated ID and timestamp
func (s *ItemService) Add(ctx context.Context, info *models.ItemInfo, content []byte) error {
	info.ID = models.ItemID(uuid.New().String())
	info.UpdatedAt = time.Now()
	return s.storage.Add(ctx, info, content)
}

// List retrieves all items for a specific user
func (s *ItemService) List(ctx context.Context, uid models.UserID) ([]models.ItemInfo, error) {
	return s.storage.ListByUser(ctx, uid)
}

// GetContent retrieves the encrypted content of an item
func (s *ItemService) GetContent(ctx context.Context, id models.ItemID) ([]byte, error) {
	return s.storage.GetContent(ctx, id)
}

// Delete removes an item by its ID
func (s *ItemService) Delete(ctx context.Context, id models.ItemID) error {
	return s.storage.DeleteItem(ctx, id)
}

// Update modifies an existing item and updates its timestamp
func (s *ItemService) Update(ctx context.Context, info *models.ItemInfo, content []byte) error {
	info.UpdatedAt = time.Now()
	return s.storage.UpdateItem(ctx, info, content)
}
