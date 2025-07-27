package card

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/rycln/gokeep/internal/shared/models"
)

const (
	nameField       = "Название"
	infoField       = "Информация"
	numberField     = "Номер"
	ownerNameField  = "Имя владельца"
	expiryDateField = "Срок действия"
	cvvField        = "CVV"
)

type Model struct {
	inputs  []textinput.Model
	focused int
}

func InitialModel() Model {
	m := Model{
		inputs: make([]textinput.Model, 6),
	}

	for i := range m.inputs {
		t := textinput.New()
		t.Width = 40
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = nameField
			t.Focus()
		case 1:
			t.Placeholder = infoField
		case 2:
			t.Placeholder = numberField
			t.CharLimit = 16
		case 3:
			t.Placeholder = ownerNameField
			t.CharLimit = 40
		case 4:
			t.Placeholder = expiryDateField
			t.CharLimit = 5
		case 5:
			t.Placeholder = cvvField
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

type (
	CardMsg struct {
		Info    *models.ItemInfo
		Content []byte
	}

	ErrMsg struct{ Err error }

	CancelMsg struct{}
)
