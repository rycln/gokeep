package update

import (
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"github.com/rycln/gokeep/client/internal/tui/items/logpass"
	"github.com/rycln/gokeep/client/internal/tui/screens/update/mocks"
	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/rycln/gokeep/client/internal/tui/shared/messages"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
)

func TestInitialModel(t *testing.T) {
	t.Run("should initialize with default values", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		timeout := 5 * time.Second

		model := InitialModel(mockService, timeout)

		assert.Equal(t, LoadState, model.state)
		assert.Equal(t, "", model.errMsg)
		assert.NotNil(t, model.logpassModel)
		assert.NotNil(t, model.cardModel)
		assert.NotNil(t, model.textModel)
		assert.NotNil(t, model.binModel)
		assert.Nil(t, model.info)
		assert.Nil(t, model.content)
		assert.Equal(t, mockService, model.service)
		assert.Equal(t, timeout, model.timeout)
	})
}

func TestSetItem(t *testing.T) {
	t.Run("should set item info and content", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)
		info := &models.ItemInfo{ID: "test-id"}
		content := []byte("test content")

		model.SetItem(info, content)
		assert.Equal(t, LoadState, model.state)
		assert.Equal(t, info, model.info)
		assert.Equal(t, content, model.content)
	})
}

func TestInit(t *testing.T) {
	t.Run("should return nil command", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)

		cmd := model.Init()
		assert.Nil(t, cmd)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("should transition to ProcessingState on ItemMsg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = UpdatePassword

		info := &models.ItemInfo{ID: "test-id"}
		content := []byte("test content")
		newModel, cmd := model.Update(messages.ItemMsg{
			Info:    info,
			Content: content,
		})

		assert.Equal(t, ProcessingState, newModel.(Model).state)
		assert.NotNil(t, cmd)
	})

	t.Run("should transition to ErrorState on ErrMsg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)

		testErr := errors.New("test error")
		newModel, _ := model.Update(messages.ErrMsg{Err: testErr})

		assert.Equal(t, ErrorState, newModel.(Model).state)
		assert.Equal(t, testErr.Error(), newModel.(Model).errMsg)
	})

	t.Run("should return CancelMsg on CancelMsg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)

		_, cmd := model.Update(messages.CancelMsg{})
		assert.NotNil(t, cmd)
	})
}

func TestHandleProcessingState(t *testing.T) {
	t.Run("should transition to ErrorState on error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = ProcessingState

		testErr := errors.New("test error")
		newModel, _ := handleProcessingState(model, ErrorMsg{testErr})

		assert.Equal(t, ErrorState, newModel.state)
		assert.Equal(t, testErr.Error(), newModel.errMsg)
	})

	t.Run("should transition to correct state on InitSuccessMsg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = ProcessingState

		newModel, _ := handleProcessingState(model, InitSuccessMsg{state: UpdatePassword})
		assert.Equal(t, UpdatePassword, newModel.state)
	})

	t.Run("should return CancelMsg on UpdateSuccessMsg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = ProcessingState

		_, cmd := handleProcessingState(model, UpdateSuccessMsg{})
		assert.NotNil(t, cmd)
	})
}

func TestHandleErrorState(t *testing.T) {
	t.Run("should return CancelMsg on Enter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = ErrorState

		_, cmd := handleErrorState(model, tea.KeyMsg{Type: tea.KeyEnter})
		assert.NotNil(t, cmd)
	})
}

func TestUpdateOperation(t *testing.T) {
	t.Run("should return UpdateSuccessMsg on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)
		model.info = &models.ItemInfo{ID: "test-id", UserID: "user-id"}
		info := &models.ItemInfo{ID: "new-id"}
		content := []byte("test content")

		mockService.EXPECT().
			Update(gomock.Any(), info, content).
			Return(nil)

		cmd := model.update(info, content)
		msg := cmd().(UpdateSuccessMsg)

		assert.NotNil(t, msg)
	})

	t.Run("should return ErrorMsg on failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)
		model.info = &models.ItemInfo{ID: "test-id", UserID: "user-id"}
		info := &models.ItemInfo{ID: "new-id"}
		content := []byte("test content")

		testErr := errors.New("update failed")
		mockService.EXPECT().
			Update(gomock.Any(), info, content).
			Return(testErr)

		cmd := model.update(info, content)
		msg := cmd().(ErrorMsg)

		assert.Equal(t, testErr, msg.Err)
	})
}

func TestInitUpdate(t *testing.T) {
	t.Run("should return error when info is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)
		model.content = []byte("test")

		cmd := model.initUpdate()
		msg := cmd().(ErrorMsg)

		assert.Equal(t, errEmptyItem, msg.Err)
	})

	t.Run("should return error when content is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)
		model.info = &models.ItemInfo{ID: "test-id"}

		cmd := model.initUpdate()
		msg := cmd().(ErrorMsg)

		assert.Equal(t, errEmptyItem, msg.Err)
	})
}

func TestView(t *testing.T) {
	t.Run("should show loading message in LoadState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = LoadState

		view := model.View()
		assert.Contains(t, view, i18n.CommonPressAnyKey)
	})

	t.Run("should show loading message in ProcessingState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = ProcessingState

		view := model.View()
		assert.Contains(t, view, i18n.CommonWait)
	})

	t.Run("should show error message in ErrorState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = ErrorState
		model.errMsg = "test error"

		view := model.View()
		assert.Contains(t, view, "test error")
	})

	t.Run("should delegate to logpass model in UpdatePassword state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemUpdater(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = UpdatePassword
		model.logpassModel = logpass.InitialModel()

		view := model.View()
		assert.NotEmpty(t, view)
	})
}
