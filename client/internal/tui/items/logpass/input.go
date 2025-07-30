package logpass

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

// Model represents a login/password entry form
type Model struct {
	input.Form // Embedded form component
}

// InitialModel creates a new login/password form with configured fields
func InitialModel() Model {
	m := Model{
		input.Form{
			Inputs: make([]textinput.Model, 4),
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
			t.Placeholder = i18n.LogPassInputLogin
		case 3:
			t.Placeholder = i18n.LogPassInputPassword
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
			ItemType: models.TypePassword,
			Name:     m.Inputs[0].Value(),
			Metadata: m.Inputs[1].Value(),
		}

		logpass := &LogPass{
			Login:    m.Inputs[2].Value(),
			Password: m.Inputs[3].Value(),
		}

		content, err := json.Marshal(logpass)
		if err != nil {
			return messages.ErrMsg{Err: fmt.Errorf("marshal error: %w", err)}
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

// SetStartData pre-populates form with existing credentials
func (m *Model) SetStartData(info *models.ItemInfo, content []byte) error {
	var logPass LogPass

	err := json.Unmarshal(content, &logPass)
	if err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	m.Inputs[0].SetValue(info.Name)
	m.Inputs[1].SetValue(info.Metadata)
	m.Inputs[2].SetValue(logPass.Login)
	m.Inputs[3].SetValue(logPass.Password)

	return nil
}
