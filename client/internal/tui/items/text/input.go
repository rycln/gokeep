package text

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

// Model represents a text entry form component
type Model struct {
	input.Form // Embedded form component
}

// InitialModel creates a new text entry form with configured fields
func InitialModel() Model {
	m := Model{
		input.Form{
			Inputs: make([]textinput.Model, 3),
		},
	}

	for i := range m.Inputs {
		t := textinput.New()
		t.Width = 30
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = i18n.InputName
			t.Focus()
		case 1:
			t.Placeholder = i18n.InputInfo
		case 2:
			t.Placeholder = i18n.TextInputContent
			t.CharLimit = 0
			t.Width = 144
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
			ItemType: models.TypeText,
			Name:     m.Inputs[0].Value(),
			Metadata: m.Inputs[1].Value(),
		}

		text := &Text{
			Content: m.Inputs[2].Value(),
		}

		content, err := json.Marshal(text)
		if err != nil {
			return messages.ErrMsg{Err: fmt.Errorf("failed to marshal text: %w", err)}
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

// SetStartData pre-populates form with existing text data
func (m *Model) SetStartData(info *models.ItemInfo, content []byte) error {
	var text Text

	err := json.Unmarshal(content, &text)
	if err != nil {
		return fmt.Errorf("failed to unmarshal text content: %w", err)
	}

	m.Inputs[0].SetValue(info.Name)
	m.Inputs[1].SetValue(info.Metadata)
	m.Inputs[2].SetValue(text.Content)

	return nil
}
