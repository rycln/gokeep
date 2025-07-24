package tui

import (
	"github.com/rycln/gokeep/internal/client/tui/auth"
	"github.com/rycln/gokeep/internal/client/tui/vault"
	"github.com/rycln/gokeep/internal/shared/models"

	tea "github.com/charmbracelet/bubbletea"
)

type rootModel struct {
	authModel  auth.Model
	vaultModel vault.Model
	current    string
	user       *models.User
}

func InitialRootModel(auth auth.Model, vault vault.Model) rootModel {
	return rootModel{
		authModel:  auth,
		vaultModel: vault,
		current:    "auth",
	}
}

func (m rootModel) Init() tea.Cmd {
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case auth.AuthSuccessMsg:
		m.user = msg.User
		m.current = "vault"
		return m, nil
	default:
		switch m.current {
		case "auth":
			updated, cmd := m.authModel.Update(msg)
			if authModel, ok := updated.(auth.Model); ok {
				m.authModel = authModel
			}
			return m, cmd
		case "vault":
			updated, cmd := m.vaultModel.Update(msg)
			if vaultModel, ok := updated.(vault.Model); ok {
				m.vaultModel = vaultModel
			}
			return m, cmd
		default:
			return m, nil
		}
	}
}

func (m rootModel) View() string {
	switch m.current {
	case "auth":
		return m.authModel.View()
	case "vault":
		return m.vaultModel.View()
	default:
		return ""
	}
}
