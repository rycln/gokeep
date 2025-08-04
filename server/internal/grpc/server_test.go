package grpc

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/gokeep/server/internal/grpc/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewGophKeeperServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockuserService(ctrl)
	mockSync := mocks.NewMocksyncService(ctrl)
	mockAuth := mocks.NewMockauthProvider(ctrl)

	t.Run("should create new server instance", func(t *testing.T) {
		server := NewGophKeeperServer(mockUser, mockSync, mockAuth, testTimeout)
		assert.NotNil(t, server)
		assert.Equal(t, mockUser, server.user)
		assert.Equal(t, mockSync, server.sync)
		assert.Equal(t, testTimeout, server.timeout)
	})
}
