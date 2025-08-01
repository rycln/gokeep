package vault

import (
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"github.com/rycln/gokeep/client/internal/tui/screens/vault/mocks"
	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
)

func TestInitialModel(t *testing.T) {
	t.Run("should initialize model with correct defaults", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		timeout := 5 * time.Second

		model := InitialModel(mockItemService, mockSyncService, timeout)

		assert.Equal(t, UpdateState, model.state)
		assert.NotNil(t, model.list)
		assert.Equal(t, mockItemService, model.itemService)
		assert.Equal(t, mockSyncService, model.syncService)
		assert.Equal(t, timeout, model.timeout)
	})
}

func TestSetUser(t *testing.T) {
	t.Run("should set user correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		user := &models.User{ID: "test-user"}

		model.SetUser(user)
		assert.Equal(t, user, model.user)
	})
}

func TestSetUpdateState(t *testing.T) {
	t.Run("should change state to UpdateState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		model.state = ListState

		model.SetUpdateState()
		assert.Equal(t, UpdateState, model.state)
	})
}

func TestInit(t *testing.T) {
	t.Run("should return error when user not set", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)

		cmd := model.Init()
		assert.NotNil(t, cmd)

		msg := cmd()
		errMsg, ok := msg.(ErrorMsg)
		assert.True(t, ok)
		assert.EqualError(t, errMsg.Err, "user not set")
	})

	t.Run("should return nil when user is set", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		model.SetUser(&models.User{ID: "test-user"})

		cmd := model.Init()
		assert.Nil(t, cmd)
	})
}

func TestLoadItems(t *testing.T) {
	t.Run("should return ItemsMsg on successful load", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		user := &models.User{ID: "test-user"}
		model.SetUser(user)

		expectedItems := []models.ItemInfo{
			{ID: "item1", Name: "Item 1"},
			{ID: "item2", Name: "Item 2"},
		}

		mockItemService.EXPECT().
			List(gomock.Any(), user.ID).
			Return(expectedItems, nil)

		cmd := model.loadItems()
		msg := cmd().(ItemsMsg)

		assert.Len(t, msg.Items, len(expectedItems))
		assert.Equal(t, expectedItems[0].ID, msg.Items[0].ID)
	})

	t.Run("should return ErrorMsg on service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		user := &models.User{ID: "test-user"}
		model.SetUser(user)

		testErr := errors.New("service error")
		mockItemService.EXPECT().
			List(gomock.Any(), user.ID).
			Return(nil, testErr)

		cmd := model.loadItems()
		msg := cmd().(ErrorMsg)

		assert.Equal(t, testErr, msg.Err)
	})
}

func TestSyncItems(t *testing.T) {
	t.Run("should return SyncSuccessMsg on successful sync", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		user := &models.User{ID: models.UserID("test-user")}
		model.SetUser(user)

		mockSyncService.EXPECT().
			SyncUserItems(gomock.Any(), user).
			Return(nil)

		cmd := model.syncItems()
		msg := cmd().(SyncSuccessMsg)

		assert.NotNil(t, msg)
	})

	t.Run("should return ErrorMsg on sync error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		user := &models.User{ID: models.UserID("test-user")}
		model.SetUser(user)

		testErr := errors.New("sync error")
		mockSyncService.EXPECT().
			SyncUserItems(gomock.Any(), user).
			Return(testErr)

		cmd := model.syncItems()
		msg := cmd().(ErrorMsg)

		assert.Equal(t, testErr, msg.Err)
	})
}

func TestDeleteItem(t *testing.T) {
	t.Run("should return DeleteSuccessMsg on successful deletion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		model.selected = &itemRender{ID: models.ItemID("test-id")}

		mockItemService.EXPECT().
			Delete(gomock.Any(), models.ItemID("test-id")).
			Return(nil)

		cmd := model.deleteItem()
		msg := cmd().(DeleteSuccessMsg)

		assert.NotNil(t, msg)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("should transition to ProcessingState from UpdateState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		model.state = UpdateState

		newModel, cmd := model.Update(nil)
		assert.Equal(t, ProcessingState, newModel.(Model).state)
		assert.NotNil(t, cmd)
	})
}

func TestHandleListState(t *testing.T) {
	t.Run("should refresh items on 'u' key", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		model.state = ListState

		newModel, cmd := handleListState(model, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}})
		assert.Equal(t, ProcessingState, newModel.state)
		assert.NotNil(t, cmd)
	})

	t.Run("should sync items on 's' key", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		model.state = ListState

		newModel, cmd := handleListState(model, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
		assert.Equal(t, ProcessingState, newModel.state)
		assert.NotNil(t, cmd)
	})
}

func TestHandleDetailState(t *testing.T) {
	t.Run("should return to list on Escape", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		model.state = DetailState
		model.selected = &itemRender{ID: "test-id"}

		newModel, _ := handleDetailState(model, tea.KeyMsg{Type: tea.KeyEsc})
		assert.Equal(t, ListState, newModel.state)
	})

	t.Run("should get content on Enter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		model.state = DetailState
		model.selected = &itemRender{ID: "test-id", ItemType: models.TypeText}

		newModel, cmd := handleDetailState(model, tea.KeyMsg{Type: tea.KeyEnter})
		assert.Equal(t, ProcessingState, newModel.state)
		assert.NotNil(t, cmd)
	})
}

func TestHandleProcessingState(t *testing.T) {
	t.Run("should handle ItemsMsg correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		model.state = ProcessingState

		testItems := []itemRender{
			{ID: "item1", Name: "Item 1"},
			{ID: "item2", Name: "Item 2"},
		}
		newModel, _ := handleProcessingState(model, ItemsMsg{Items: testItems})

		assert.Equal(t, ListState, newModel.state)
		assert.Equal(t, testItems, newModel.items)
	})

	t.Run("should handle ErrorMsg correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		model.state = ProcessingState

		testErr := errors.New("test error")
		newModel, _ := handleProcessingState(model, ErrorMsg{Err: testErr})

		assert.Equal(t, ErrorState, newModel.state)
		assert.Equal(t, testErr.Error(), newModel.errMsg)
	})

	t.Run("should handle SyncSuccessMsg correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		model.state = ProcessingState

		newModel, _ := handleProcessingState(model, SyncSuccessMsg{})
		assert.Equal(t, UpdateState, newModel.state)
	})
}

func TestView(t *testing.T) {
	t.Run("should show loading message in ProcessingState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		model.state = ProcessingState

		view := model.View()
		assert.Contains(t, view, i18n.CommonWait)
	})

	t.Run("should show error message in ErrorState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockItemService := mocks.NewMockitemService(ctrl)
		mockSyncService := mocks.NewMocksyncService(ctrl)
		model := InitialModel(mockItemService, mockSyncService, time.Second)
		model.state = ErrorState
		model.errMsg = "test error"

		view := model.View()
		assert.Contains(t, view, "test error")
	})
}
