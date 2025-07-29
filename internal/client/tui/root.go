package tui

import (
	"github.com/rycln/gokeep/internal/client/tui/screens/add"
	"github.com/rycln/gokeep/internal/client/tui/screens/auth"
	"github.com/rycln/gokeep/internal/client/tui/screens/vault"

	tea "github.com/charmbracelet/bubbletea"
)

type model int

const (
	AuthModel model = iota
	VaultModel
	AddModel
)

type rootModel struct {
	authModel  auth.Model
	vaultModel vault.Model
	addModel   add.Model
	current    model
}

func InitialRootModel(auth auth.Model, vault vault.Model, add add.Model) rootModel {
	return rootModel{
		authModel:  auth,
		vaultModel: vault,
		addModel:   add,
		current:    AuthModel,
	}
}

func (m rootModel) Init() tea.Cmd {
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case auth.AuthSuccessMsg:
		m.vaultModel.SetUser(msg.User)
		m.current = VaultModel
		return m, nil
	case vault.AddItemReqMsg:
		m.addModel.SetUser(msg.User)
		m.current = AddModel
		return m, nil
	case add.CancelMsg:
		m.vaultModel.SetUpdateState()
		m.current = VaultModel
		return m, nil
	default:
		switch m.current {
		case AuthModel:
			updated, cmd := m.authModel.Update(msg)
			if authModel, ok := updated.(auth.Model); ok {
				m.authModel = authModel
			}
			return m, cmd
		case VaultModel:
			updated, cmd := m.vaultModel.Update(msg)
			if vaultModel, ok := updated.(vault.Model); ok {
				m.vaultModel = vaultModel
			}
			return m, cmd
		case AddModel:
			updated, cmd := m.addModel.Update(msg)
			if addModel, ok := updated.(add.Model); ok {
				m.addModel = addModel
			}
			return m, cmd
		default:
			return m, nil
		}
	}
}

func (m rootModel) View() string {
	switch m.current {
	case AuthModel:
		return m.authModel.View()
	case VaultModel:
		return m.vaultModel.View()
	case AddModel:
		return m.addModel.View()
	default:
		return ""
	}
}
