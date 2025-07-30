package card

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/rycln/gokeep/client/internal/tui/shared/input"
	"github.com/rycln/gokeep/client/internal/tui/shared/messages"
	"github.com/rycln/gokeep/shared/models"
)

// Model represents a credit card entry form
type Model struct {
	input.Form // Embedded form component
}

// InitialModel creates a new card entry form with configured fields
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
			t.Placeholder = i18n.InputName
			t.Focus()
		case 1:
			t.Placeholder = i18n.InputInfo
		case 2:
			t.Placeholder = i18n.CardInputNumber
			t.CharLimit = 16
		case 3:
			t.Placeholder = i18n.CardInputHolderName
			t.CharLimit = 40
		case 4:
			t.Placeholder = i18n.CardInputExpiryDate
			t.CharLimit = 5
		case 5:
			t.Placeholder = i18n.CardInputCVV
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.Inputs[i] = t
	}

	return m
}

// Update handles form input and submission events
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

// send packages form data into a message for processing
func (m Model) send() tea.Cmd {
	return func() tea.Msg {
		info := &models.ItemInfo{
			ItemType: models.TypeCard,
			Name:     m.Inputs[0].Value(),
			Metadata: m.Inputs[1].Value(),
		}

		card := &Card{
			CardNumber: m.Inputs[2].Value(),
			CardOwner:  m.Inputs[3].Value(),
			ExpiryDate: m.Inputs[4].Value(),
			CVV:        m.Inputs[5].Value(),
		}

		content, err := json.Marshal(card)
		if err != nil {
			return messages.ErrMsg{Err: err}
		}

		for i := range m.Inputs {
			m.Inputs[i].Reset()
		}

		return messages.ItemMsg{
			Info:    info,
			Content: content,
		}
	}
}

// SetStartData pre-populates form with existing card data
func (m *Model) SetStartData(info *models.ItemInfo, content []byte) error {
	var card Card

	err := json.Unmarshal(content, &card)
	if err != nil {
		return fmt.Errorf("failed to unmarshal card data: %w", err)
	}

	m.Inputs[0].SetValue(info.Name)
	m.Inputs[1].SetValue(info.Metadata)
	m.Inputs[2].SetValue(card.CardNumber)
	m.Inputs[3].SetValue(card.CardOwner)
	m.Inputs[4].SetValue(card.ExpiryDate)
	m.Inputs[5].SetValue(card.CVV)

	return nil
}
