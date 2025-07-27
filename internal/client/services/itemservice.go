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
	GetContent(context.Context, string) ([]byte, error)
}

type itemStorage interface {
	itemStorer
	itemGetter
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
func (s *ItemService) GetContent(ctx context.Context, name string) ([]byte, error) {
	return s.storage.GetContent(ctx, name)
}
