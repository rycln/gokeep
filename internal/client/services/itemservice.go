package services

import (
	"context"

	"github.com/rycln/gokeep/internal/shared/models"
)

type itemStorer interface {
	AddItem(context.Context, *models.ItemInfo, []byte) error
}

type itemGetter interface {
	GetUserItemsInfo(context.Context, models.UserID) ([]models.ItemInfo, error)
	GetContentByName(context.Context, string) ([]byte, error)
}

type itemStorage interface {
	itemStorer
	itemGetter
}

type ItemService struct {
	strg itemStorage
}

func NewItemService(strg itemStorage) *ItemService {
	return &ItemService{
		strg: strg,
	}
}
