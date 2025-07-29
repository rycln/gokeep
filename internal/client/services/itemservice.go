package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rycln/gokeep/internal/shared/models"
)

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

type itemStorage interface {
	itemStorer
	itemGetter
	itemDeleter
	itemUpdater
}

type ItemService struct {
	storage itemStorage
}

func NewItemService(storage itemStorage) *ItemService {
	return &ItemService{
		storage: storage,
	}
}

// добавить шифрование
func (s *ItemService) Add(ctx context.Context, info *models.ItemInfo, content []byte) error {
	info.ID = models.ItemID(uuid.New().String())
	info.UpdatedAt = time.Now()
	return s.storage.Add(ctx, info, content)
}

func (s *ItemService) List(ctx context.Context, uid models.UserID) ([]models.ItemInfo, error) {
	return s.storage.ListByUser(ctx, uid)
}

// добавить дешифровку
func (s *ItemService) GetContent(ctx context.Context, id models.ItemID) ([]byte, error) {
	return s.storage.GetContent(ctx, id)
}

func (s *ItemService) Delete(ctx context.Context, id models.ItemID) error {
	return s.storage.DeleteItem(ctx, id)
}

func (s *ItemService) Update(ctx context.Context, info *models.ItemInfo, content []byte) error {
	info.UpdatedAt = time.Now()
	return s.storage.UpdateItem(ctx, info, content)
}
