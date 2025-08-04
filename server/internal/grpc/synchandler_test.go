package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	pb "github.com/rycln/gokeep/pkg/gen/grpc/gophkeeper"
	"github.com/rycln/gokeep/server/internal/grpc/mocks"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGophKeeperServer_Sync(t *testing.T) {
	now := time.Now().UTC()
	testTimeout := 5 * time.Second

	t.Run("successful sync with items", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUser := mocks.NewMockuserService(ctrl)
		mockSync := mocks.NewMocksyncService(ctrl)
		mockAuth := mocks.NewMockauthProvider(ctrl)
		handler := NewGophKeeperServer(mockUser, mockSync, mockAuth, testTimeout)

		req := &pb.SyncRequest{
			Items: []*pb.Item{
				{
					Id:        "item1",
					UserId:    "user1",
					Type:      "type1",
					Name:      "name1",
					Metadata:  "meta1",
					Data:      []byte("data1"),
					UpdatedAt: timestamppb.New(now),
					IsDeleted: false,
				},
			},
		}

		expectedItems := []models.Item{
			{
				ID:        "item1",
				UserID:    "user1",
				ItemType:  "type1",
				Name:      "name1",
				Metadata:  "meta1",
				Data:      []byte("data1"),
				UpdatedAt: now,
				IsDeleted: false,
			},
		}

		returnItems := []models.Item{
			{
				ID:        "server-item1",
				UserID:    "user1",
				ItemType:  "type1",
				Name:      "server-name1",
				Metadata:  "server-meta1",
				Data:      []byte("server-data1"),
				UpdatedAt: now.Add(time.Hour),
				IsDeleted: true,
			},
		}

		mockSync.EXPECT().
			SyncItems(gomock.Any(), gomock.Eq(expectedItems)).
			Return(returnItems, nil)

		resp, err := handler.Sync(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Items, 1)

		assert.Equal(t, "server-item1", resp.Items[0].Id)
		assert.Equal(t, "user1", resp.Items[0].UserId)
		assert.Equal(t, "type1", resp.Items[0].Type)
		assert.Equal(t, "server-name1", resp.Items[0].Name)
		assert.Equal(t, "server-meta1", resp.Items[0].Metadata)
		assert.Equal(t, []byte("server-data1"), resp.Items[0].Data)
		assert.True(t, resp.Items[0].IsDeleted)
		assert.NotNil(t, resp.Items[0].UpdatedAt)
	})

	t.Run("empty request", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUser := mocks.NewMockuserService(ctrl)
		mockSync := mocks.NewMocksyncService(ctrl)
		mockAuth := mocks.NewMockauthProvider(ctrl)
		handler := NewGophKeeperServer(mockUser, mockSync, mockAuth, testTimeout)

		req := &pb.SyncRequest{Items: []*pb.Item{}}

		mockSync.EXPECT().
			SyncItems(gomock.Any(), gomock.Eq([]models.Item{})).
			Return([]models.Item{}, nil)

		resp, err := handler.Sync(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Empty(t, resp.Items)
	})

	t.Run("service returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUser := mocks.NewMockuserService(ctrl)
		mockSync := mocks.NewMocksyncService(ctrl)
		mockAuth := mocks.NewMockauthProvider(ctrl)
		handler := NewGophKeeperServer(mockUser, mockSync, mockAuth, testTimeout)

		req := &pb.SyncRequest{
			Items: []*pb.Item{{Id: "item1"}},
		}

		expectedErr := errors.New("sync failed")
		mockSync.EXPECT().
			SyncItems(gomock.Any(), gomock.Any()).
			Return(nil, expectedErr)

		resp, err := handler.Sync(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Contains(t, err.Error(), expectedErr.Error())
	})
}
