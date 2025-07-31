package add

import (
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"github.com/rycln/gokeep/client/internal/tui/items/logpass"
	"github.com/rycln/gokeep/client/internal/tui/screens/add/mocks"
	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
)

func TestInitialModel(t *testing.T) {
	t.Run("should initialize with default values", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		timeout := 5 * time.Second

		model := InitialModel(mockService, timeout)

		assert.Equal(t, SelectState, model.state)
		assert.Len(t, model.choices, 4)
		assert.Equal(t, 0, model.cursor)
		assert.Equal(t, "", model.errMsg)
		assert.NotNil(t, model.logpassModel)
		assert.NotNil(t, model.cardModel)
		assert.NotNil(t, model.textModel)
		assert.NotNil(t, model.binModel)
		assert.Equal(t, mockService, model.service)
		assert.Equal(t, timeout, model.timeout)
	})
}

func TestSetUser(t *testing.T) {
	t.Run("should set user correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)
		user := &models.User{ID: "test-user"}

		model.SetUser(user)
		assert.Equal(t, user, model.user)
	})
}

func TestInit(t *testing.T) {
	t.Run("should return nil command", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)

		cmd := model.Init()
		assert.Nil(t, cmd)
	})
}

func TestHandleSelectState(t *testing.T) {
	t.Run("should cancel on Esc key", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)

		_, cmd := handleSelectState(model, tea.KeyMsg{Type: tea.KeyEsc})
		assert.NotNil(t, cmd)
	})

	t.Run("should move cursor up", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)
		model.cursor = 1

		newModel, _ := handleSelectState(model, tea.KeyMsg{Type: tea.KeyUp})
		assert.Equal(t, 0, newModel.cursor)
	})

	t.Run("should move cursor down", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)

		newModel, _ := handleSelectState(model, tea.KeyMsg{Type: tea.KeyDown})
		assert.Equal(t, 1, newModel.cursor)
	})

	t.Run("should select password type on Enter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)
		model.cursor = 0 // Password is first choice

		newModel, _ := handleSelectState(model, tea.KeyMsg{Type: tea.KeyEnter})
		assert.Equal(t, AddPassword, newModel.state)
	})

	t.Run("should select card type on Enter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)
		model.cursor = 1 // Card is second choice

		newModel, _ := handleSelectState(model, tea.KeyMsg{Type: tea.KeyEnter})
		assert.Equal(t, AddCard, newModel.state)
	})
}

func TestHandleProcessingState(t *testing.T) {
	t.Run("should transition to ErrorState on error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = ProcessingState

		testErr := errors.New("test error")
		newModel, _ := handleProcessingState(model, ErrorMsg{testErr})

		assert.Equal(t, ErrorState, newModel.state)
		assert.Equal(t, testErr.Error(), newModel.errMsg)
	})

	t.Run("should return to SelectState on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = ProcessingState

		newModel, _ := handleProcessingState(model, AddSuccessMsg{})
		assert.Equal(t, SelectState, newModel.state)
	})
}

func TestHandleErrorState(t *testing.T) {
	t.Run("should return to SelectState on Enter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = ErrorState

		newModel, _ := handleErrorState(model, tea.KeyMsg{Type: tea.KeyEnter})
		assert.Equal(t, SelectState, newModel.state)
	})
}

func TestAdd(t *testing.T) {
	t.Run("should return AddSuccessMsg on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)
		info := &models.ItemInfo{ID: "test-id"}
		content := []byte("test content")

		mockService.EXPECT().
			Add(gomock.Any(), info, content).
			Return(nil)

		cmd := model.add(info, content)
		msg := cmd().(AddSuccessMsg)

		assert.NotNil(t, msg)
	})

	t.Run("should return ErrorMsg on failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)
		info := &models.ItemInfo{ID: "test-id"}
		content := []byte("test content")

		testErr := errors.New("add failed")
		mockService.EXPECT().
			Add(gomock.Any(), info, content).
			Return(testErr)

		cmd := model.add(info, content)
		msg := cmd().(ErrorMsg)

		assert.Equal(t, testErr, msg.Err)
	})
}

func TestView(t *testing.T) {
	t.Run("should show selection view in SelectState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = SelectState

		view := model.View()
		assert.Contains(t, view, i18n.AddSelectPrompt)
	})

	t.Run("should show loading message in ProcessingState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = ProcessingState

		view := model.View()
		assert.Contains(t, view, i18n.CommonWait)
	})

	t.Run("should show error message in ErrorState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = ErrorState
		model.errMsg = "test error"

		view := model.View()
		assert.Contains(t, view, "test error")
	})

	t.Run("should delegate to logpass model in AddPassword state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = AddPassword
		model.logpassModel = logpass.InitialModel() // Можно замокать при необходимости

		view := model.View()
		assert.NotEmpty(t, view)
	})
}

func TestHandleItemModelUpdate(t *testing.T) {
	t.Run("should update logpass model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockitemAdder(ctrl)
		model := InitialModel(mockService, time.Second)
		model.state = AddPassword

		mockLogpass := logpass.InitialModel() // Можно замокать при необходимости
		model.logpassModel = mockLogpass

		newModel, _ := handleItemModelUpdate(model, nil, &model.logpassModel)
		assert.Equal(t, mockLogpass, newModel.logpassModel)
	})
}
