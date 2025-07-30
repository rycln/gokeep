package server

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/gokeep/server/internal/grpc/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewGophKeeperServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockuserService(ctrl)

	t.Run("should create new server instance", func(t *testing.T) {
		server := NewGophKeeperServer(mockService, testTimeout)
		assert.NotNil(t, server)
		assert.Equal(t, mockService, server.uservice)
		assert.Equal(t, testTimeout, server.timeout)
	})
}
