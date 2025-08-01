package services

import (
	"context"

	"github.com/rycln/gokeep/shared/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// syncAPI defines the interface for synchronization operations with remote server
type syncAPI interface {
	// Sync performs bidirectional synchronization between client and server items
	Sync(context.Context, []models.Item, string) ([]models.Item, error)
}

// syncStorage defines the interface for item storage operations
type syncStorage interface {
	// GetAllUserItems retrieves all items for specified user from local storage
	GetAllUserItems(context.Context, models.UserID) ([]models.Item, error)
	// ReplaceAllUserItems completely replaces user's items in local storage
	ReplaceAllUserItems(context.Context, models.UserID, []models.Item) error
}

// SyncService handles synchronization between local storage and remote server
type SyncService struct {
	sync syncAPI     // Remote synchronization API
	strg syncStorage // Local items storage
}

// NewSyncService creates a new SyncService instance
func NewSyncService(sync syncAPI, strg syncStorage) *SyncService {
	return &SyncService{
		sync: sync,
		strg: strg,
	}
}

// SyncUserItems performs full synchronization cycle for user's items:
func (s *SyncService) SyncUserItems(ctx context.Context, user *models.User) error {
	clientItems, err := s.strg.GetAllUserItems(ctx, user.ID)
	if err != nil {
		return err
	}

	serverItems, err := s.sync.Sync(ctx, clientItems, user.JWT)
	if err != nil {
		return err
	}

	err = s.strg.ReplaceAllUserItems(ctx, user.ID, serverItems)
	if err != nil {
		return err
	}

	return nil
}
