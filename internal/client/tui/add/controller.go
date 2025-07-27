package add

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/internal/client/tui/add/logpass"
	"github.com/rycln/gokeep/internal/shared/models"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case logpass.LogPassMsg:
		info := msg.Info
		info.UserID = m.user.ID
		m.state = ProcessingState
		return m, m.add(info, msg.Content)
	case logpass.ErrMsg:
		m.errMsg = msg.Err.Error()
		m.state = ErrorState
		return m, nil
	case logpass.CancelMsg:
		m.state = SelectState
		return m, nil
	default:
		switch m.state {
		case AddPassword:
			updated, cmd := m.logpassModel.Update(msg)
			if logpassModel, ok := updated.(logpass.Model); ok {
				m.logpassModel = logpassModel
			}
			return m, cmd
		default:
			switch msg := msg.(type) {
			case tea.KeyMsg:
				if msg.Type == tea.KeyCtrlC {
					return m, tea.Quit
				}
				return m.handleKeyMsg(msg)

			case ErrorMsg:
				m.errMsg = msg.Err.Error()
				m.state = ErrorState
				return m, nil

			case AddSuccessMsg:
				m.state = SelectState
				return m, nil

			default:
				return m, nil
			}
		}
	}
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch m.state {
	case SelectState:
		switch msg.Type {
		case tea.KeyEsc:
			return m, func() tea.Msg { return CancelMsg{} }
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case tea.KeyEnter:
			choice := m.choices[m.cursor]
			switch choice {
			case password:
				m.state = AddPassword
			}
		}

	case ErrorState:
		if msg.Type == tea.KeyEnter {
			m.state = SelectState
		}
	}

	return m, nil
}

func (m Model) add(info *models.ItemInfo, content []byte) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
		defer cancel()

		err := m.service.Add(ctx, info, content)
		if err != nil {
			return ErrorMsg{err}
		}

		return AddSuccessMsg{}
	}
}
