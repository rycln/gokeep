package card

import (
	"encoding/json"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/internal/client/tui/input"
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
	input.Form
}

func InitialModel() Model {
	m := Model{
		input.Form{
			Inputs: make([]textinput.Model, 6),
		},
	}

	for i := range m.Inputs {
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

		m.Inputs[i] = t
	}

	return m
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, m.send()
		}
	}

	updated, cmd := m.Form.Update(msg)
	if inputForm, ok := updated.(input.Form); ok {
		m.Form = inputForm
	}
	return m, cmd
}

func (m Model) send() tea.Cmd {
	return func() tea.Msg {
		info := &models.ItemInfo{
			ItemType: models.TypeCard,
			Name:     m.Inputs[0].Value(),
			Metadata: m.Inputs[1].Value(),
		}

		card := &models.Card{
			CardNumber: m.Inputs[2].Value(),
			CardOwner:  m.Inputs[3].Value(),
			ExpiryDate: m.Inputs[4].Value(),
			CVV:        m.Inputs[5].Value(),
		}

		content, err := json.Marshal(card)
		if err != nil {
			return input.ErrMsg{Err: err}
		}

		return input.ItemMsg{
			Info:    info,
			Content: content,
		}
	}
}
