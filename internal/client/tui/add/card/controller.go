package card

import (
	"encoding/json"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rycln/gokeep/internal/shared/models"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle      = lipgloss.NewStyle()
)

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			return m, func() tea.Msg { return CancelMsg{} }
		case tea.KeyEnter:
			return m, m.send()
		case tea.KeyUp:
			m.focused--
			if m.focused < 0 {
				m.focused = len(m.inputs) - 1
			}
			cmd = m.updateFocus()
			return m, cmd
		case tea.KeyDown:
			m.focused++
			if m.focused >= len(m.inputs) {
				m.focused = 0
			}
			cmd = m.updateFocus()
			return m, cmd
		}
	}

	cmd = m.updateInputs(msg)
	return m, cmd
}

func (m *Model) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		if i == m.focused {
			cmds[i] = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
		} else {
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = noStyle
			m.inputs[i].TextStyle = noStyle
		}
	}
	return tea.Batch(cmds...)
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m Model) send() tea.Cmd {
	return func() tea.Msg {
		info := &models.ItemInfo{
			ItemType: models.TypeCard,
			Name:     m.inputs[0].Value(),
			Metadata: m.inputs[1].Value(),
		}

		card := &models.Card{
			CardNumber: m.inputs[2].Value(),
			CardOwner:  m.inputs[3].Value(),
			ExpiryDate: m.inputs[4].Value(),
			CVV:        m.inputs[5].Value(),
		}

		content, err := json.Marshal(card)
		if err != nil {
			return ErrMsg{Err: err}
		}

		return CardMsg{
			Info:    info,
			Content: content,
		}
	}
}
