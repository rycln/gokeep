package input

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rycln/gokeep/client/internal/tui/shared/messages"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle      = lipgloss.NewStyle()
)

func (m Form) Init() tea.Cmd {
	return textinput.Blink
}

func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			return m, func() tea.Msg { return messages.CancelMsg{} }
		case tea.KeyUp:
			m.Focused--
			if m.Focused < 0 {
				m.Focused = len(m.Inputs) - 1
			}
			cmd = m.updateFocus()
			return m, cmd
		case tea.KeyDown:
			m.Focused++
			if m.Focused >= len(m.Inputs) {
				m.Focused = 0
			}
			cmd = m.updateFocus()
			return m, cmd
		}
	}

	cmd = m.updateInputs(msg)
	return m, cmd
}

func (m *Form) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))
	for i := range m.Inputs {
		if i == m.Focused {
			cmds[i] = m.Inputs[i].Focus()
			m.Inputs[i].PromptStyle = focusedStyle
			m.Inputs[i].TextStyle = focusedStyle
		} else {
			m.Inputs[i].Blur()
			m.Inputs[i].PromptStyle = noStyle
			m.Inputs[i].TextStyle = noStyle
		}
	}
	return tea.Batch(cmds...)
}

func (m *Form) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}
