package logpass

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/rycln/gokeep/internal/shared/models"
)

const (
	nameField     = "Название"
	infoField     = "Информация"
	loginField    = "Логин"
	passwordField = "Пароль"
)

type Model struct {
	inputs  []textinput.Model
	focused int
}

func InitialModel() Model {
	m := Model{
		inputs: make([]textinput.Model, 4),
	}

	for i := range m.inputs {
		t := textinput.New()
		t.Width = 30
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = nameField
			t.Focus()
		case 1:
			t.Placeholder = infoField
		case 2:
			t.Placeholder = loginField
		case 3:
			t.Placeholder = passwordField
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

type (
	LogPassMsg struct {
		Info    *models.ItemInfo
		Content []byte
	}

	ErrMsg struct{ Err error }

	CancelMsg struct{}
)
