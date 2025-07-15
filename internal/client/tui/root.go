package tui

import (
	"github.com/rycln/gokeep/internal/client/tui/auth"

	tea "github.com/charmbracelet/bubbletea"
)

type rootModel struct {
	authModel auth.Model
	current   string
}

func InitialRootModel(auth auth.Model) rootModel {
	return rootModel{
		authModel: auth,
		current:   "auth",
	}
}

func (m rootModel) Init() tea.Cmd {
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.current {
	case "auth":
		updated, cmd := m.authModel.Update(msg)
		if authModel, ok := updated.(auth.Model); ok {
			m.authModel = authModel
			if authModel.GetState() == auth.SuccessState {
				//переход к хранилищу
			}
		}
		return m, cmd
	default:
		return m, nil
	}
}

func (m rootModel) View() string {
	switch m.current {
	case "auth":
		return m.authModel.View()
	default:
		return ""
	}
}
