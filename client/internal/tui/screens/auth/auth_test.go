package auth

import (
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"github.com/rycln/gokeep/client/internal/tui/screens/auth/mocks"
	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
)

func TestInitialModel(t *testing.T) {
	t.Run("should initialize with default values", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		timeout := 5 * time.Second

		model := InitialModel(mockService, mockKey, timeout)

		assert.Equal(t, LoginState, model.state)
		assert.Equal(t, UsernameField, model.activeField)
		assert.Equal(t, "", model.username)
		assert.Equal(t, "", model.password)
		assert.Equal(t, "", model.errMsg)
		assert.Equal(t, mockService, model.service)
		assert.Equal(t, timeout, model.timeout)
	})
}

func TestGetState(t *testing.T) {
	t.Run("should return current state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.state = RegisterState

		assert.Equal(t, RegisterState, model.GetState())
	})
}

func TestInit(t *testing.T) {
	t.Run("should return nil command", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)

		cmd := model.Init()
		assert.Nil(t, cmd)
	})
}

func TestHandleAuthInput(t *testing.T) {
	t.Run("should switch to processing state on Enter in login state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.state = LoginState

		newModel, cmd := handleAuthInput(model, tea.KeyMsg{Type: tea.KeyEnter})
		assert.Equal(t, ProcessingState, newModel.state)
		assert.NotNil(t, cmd)
	})

	t.Run("should switch to processing state on Enter in register state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.state = RegisterState

		newModel, cmd := handleAuthInput(model, tea.KeyMsg{Type: tea.KeyEnter})
		assert.Equal(t, ProcessingState, newModel.state)
		assert.NotNil(t, cmd)
	})

	t.Run("should toggle between login and register states on Tab", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)

		newModel, _ := handleAuthInput(model, tea.KeyMsg{Type: tea.KeyTab})
		assert.Equal(t, RegisterState, newModel.state)

		newModel, _ = handleAuthInput(newModel, tea.KeyMsg{Type: tea.KeyTab})
		assert.Equal(t, LoginState, newModel.state)
	})

	t.Run("should toggle between username and password fields on Up/Down", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)

		newModel, _ := handleAuthInput(model, tea.KeyMsg{Type: tea.KeyDown})
		assert.Equal(t, PasswordField, newModel.activeField)

		newModel, _ = handleAuthInput(newModel, tea.KeyMsg{Type: tea.KeyUp})
		assert.Equal(t, UsernameField, newModel.activeField)
	})

	t.Run("should update username when typing in username field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.activeField = UsernameField

		newModel, _ := handleAuthInput(model, tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'t', 'e', 's', 't'},
		})
		assert.Equal(t, "test", newModel.username)
	})

	t.Run("should update password when typing in password field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.activeField = PasswordField

		newModel, _ := handleAuthInput(model, tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'p', 'a', 's', 's'},
		})
		assert.Equal(t, "pass", newModel.password)
	})

	t.Run("should ignore space key", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.activeField = UsernameField

		newModel, _ := handleAuthInput(model, tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{' '},
		})
		assert.Equal(t, "", newModel.username)
	})

	t.Run("should handle backspace in username field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.activeField = UsernameField
		model.username = "test"

		newModel, _ := handleAuthInput(model, tea.KeyMsg{Type: tea.KeyBackspace})
		assert.Equal(t, "tes", newModel.username)
	})

	t.Run("should handle backspace in password field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.activeField = PasswordField
		model.password = "pass"

		newModel, _ := handleAuthInput(model, tea.KeyMsg{Type: tea.KeyBackspace})
		assert.Equal(t, "pas", newModel.password)
	})
}

func TestLogin(t *testing.T) {
	t.Run("should return AuthSuccessMsg on successful login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.username = "testuser"
		model.password = "testpass"

		expectedUser := &models.User{ID: "user123"}
		mockService.EXPECT().
			UserLogin(gomock.Any(), &models.UserLoginReq{
				Username: "testuser",
				Password: "testpass",
			}).
			Return(expectedUser, nil)

		cmd := model.login()
		msg := cmd().(AuthSuccessMsg)

		assert.Equal(t, expectedUser, msg.User)
	})

	t.Run("should return LoginErrorMsg on failed login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.username = "testuser"
		model.password = "testpass"

		testErr := errors.New("login failed")
		mockService.EXPECT().
			UserLogin(gomock.Any(), gomock.Any()).
			Return(nil, testErr)

		cmd := model.login()
		msg := cmd().(LoginErrorMsg)

		assert.Equal(t, testErr, msg.Err)
	})
}

func TestRegister(t *testing.T) {
	t.Run("should return AuthSuccessMsg on successful registration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.username = "newuser"
		model.password = "newpass"
		testSalt := "salt"

		expectedUser := &models.User{ID: "user456"}

		mockKey.EXPECT().GenerateSalt().Return([]byte(testSalt), nil)
		mockService.EXPECT().
			UserRegister(gomock.Any(), &models.UserRegReq{
				Username: "newuser",
				Password: "newpass",
				Salt:     testSalt,
			}).
			Return(expectedUser, nil)

		cmd := model.register()
		msg := cmd().(AuthSuccessMsg)

		assert.Equal(t, expectedUser, msg.User)
	})

	t.Run("should return RegisterErrorMsg on failed registration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.username = "newuser"
		model.password = "newpass"
		testSalt := "salt"

		testErr := errors.New("registration failed")
		mockKey.EXPECT().GenerateSalt().Return([]byte(testSalt), nil)
		mockService.EXPECT().
			UserRegister(gomock.Any(), gomock.Any()).
			Return(nil, testErr)

		cmd := model.register()
		msg := cmd().(RegisterErrorMsg)

		assert.Equal(t, testErr, msg.Err)
	})
}

func TestHandleProcessingState(t *testing.T) {
	t.Run("should transition to ErrorState on LoginErrorMsg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.state = ProcessingState

		testErr := errors.New("login error")
		newModel, _ := handleProcessingState(model, LoginErrorMsg{testErr})

		assert.Equal(t, ErrorState, newModel.state)
		assert.Equal(t, testErr.Error(), newModel.errMsg)
	})

	t.Run("should transition to ErrorState on RegisterErrorMsg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.state = ProcessingState

		testErr := errors.New("register error")
		newModel, _ := handleProcessingState(model, RegisterErrorMsg{testErr})

		assert.Equal(t, ErrorState, newModel.state)
		assert.Equal(t, testErr.Error(), newModel.errMsg)
	})

	t.Run("should return AuthSuccessMsg unchanged", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.state = ProcessingState

		user := &models.User{ID: "testuser"}
		newModel, cmd := handleProcessingState(model, AuthSuccessMsg{user})

		assert.Equal(t, model, newModel)
		assert.NotNil(t, cmd)
	})
}

func TestHandleErrorState(t *testing.T) {
	t.Run("should transition to LoginState on Enter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.state = ErrorState

		newModel, _ := handleErrorState(model, tea.KeyMsg{Type: tea.KeyEnter})
		assert.Equal(t, LoginState, newModel.state)
	})

	t.Run("should not change state on other keys", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.state = ErrorState

		newModel, _ := handleErrorState(model, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
		assert.Equal(t, ErrorState, newModel.state)
	})
}

func TestView(t *testing.T) {
	t.Run("should show loading message in ProcessingState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.state = ProcessingState

		view := model.View()
		assert.Contains(t, view, i18n.CommonWait)
	})

	t.Run("should show error message in ErrorState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.state = ErrorState
		model.errMsg = "test error"

		view := model.View()
		assert.Contains(t, view, "test error")
	})

	t.Run("should render login form in LoginState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.state = LoginState
		model.username = "user"
		model.password = "pass"

		view := model.View()
		assert.Contains(t, view, i18n.AuthLoginTitle)
		assert.Contains(t, view, "user")
		assert.Contains(t, view, "••••")
	})

	t.Run("should render register form in RegisterState", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKey := mocks.NewMockkeyProvider(ctrl)
		mockService := mocks.NewMockauthService(ctrl)
		model := InitialModel(mockService, mockKey, time.Second)
		model.state = RegisterState
		model.username = "newuser"
		model.password = "newpass"

		view := model.View()
		assert.Contains(t, view, i18n.AuthRegisterTitle)
		assert.Contains(t, view, "newuser")
		assert.Contains(t, view, "••••••")
	})
}
