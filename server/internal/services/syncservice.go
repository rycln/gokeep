package services

import (
	"context"

	"github.com/rycln/gokeep/shared/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// itemFetcher defines interface for fetching user items.
type itemFetcher interface {
	GetUserItems(context.Context, models.UserID) ([]models.Item, error)
}

// itemAdder defines interface for adding new items.
type itemAdder interface {
	AddItem(context.Context, *models.Item) error
}

// itemDeleter defines interface for deleting items.
type itemDeleter interface {
	DeleteItem(context.Context, models.ItemID, models.UserID) error
}

// itemStorage combines all item-related storage operations.
type itemStorage interface {
	itemFetcher
	itemAdder
	itemDeleter
}

// uidFetcher defines interface for getting user ID from context.
type uidFetcher interface {
	GetUserIDFromCtx(context.Context) (models.UserID, error)
}

// SyncService handles item synchronization operations.
type SyncService struct {
	strg itemStorage
	auth uidFetcher
}

// NewSyncService creates a new SyncService instance.
func NewSyncService(strg itemStorage, auth uidFetcher) *SyncService {
	return &SyncService{
		strg: strg,
		auth: auth,
	}
}

// SyncItems synchronizes user items between client and server.
func (s *SyncService) SyncItems(ctx context.Context, reqitems []models.Item) ([]models.Item, error) {
	uid, err := s.auth.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	for _, item := range reqitems {
		if item.IsDeleted {
			err := s.strg.DeleteItem(ctx, item.ID, uid)
			if err != nil {
				return nil, err
			}
		} else {
			err := s.strg.AddItem(ctx, &item)
			if err != nil {
				return nil, err
			}
		}
	}

	resitems, err := s.strg.GetUserItems(ctx, uid)
	if err != nil {
		return nil, err
	}

	return resitems, nil
}
