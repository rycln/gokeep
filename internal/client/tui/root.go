package tui

import (
	"github.com/rycln/gokeep/internal/client/tui/screens/add"
	"github.com/rycln/gokeep/internal/client/tui/screens/auth"
	"github.com/rycln/gokeep/internal/client/tui/screens/update"
	"github.com/rycln/gokeep/internal/client/tui/screens/vault"

	tea "github.com/charmbracelet/bubbletea"
)

type model int

const (
	AuthModel model = iota
	VaultModel
	AddModel
	UpdateModel
)

type rootModel struct {
	authModel   auth.Model
	vaultModel  vault.Model
	addModel    add.Model
	updateModel update.Model
	current     model
}

func InitialRootModel(auth auth.Model, vault vault.Model, add add.Model, update update.Model) rootModel {
	return rootModel{
		authModel:   auth,
		vaultModel:  vault,
		addModel:    add,
		updateModel: update,
		current:     AuthModel,
	}
}

func (m rootModel) Init() tea.Cmd {
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case add.CancelMsg, update.CancelMsg:
		m.vaultModel.SetUpdateState()
		m.current = VaultModel
		return m, nil
	default:
		switch m.current {
		case AuthModel:
			return handleAuthModel(m, msg)
		case VaultModel:
			return handleVaultModel(m, msg)
		case AddModel:
			return handleAddModel(m, msg)
		case UpdateModel:
			return handleUpdateModel(m, msg)
		default:
			return m, nil
		}
	}
}

func handleAuthModel(m rootModel, msg tea.Msg) (rootModel, tea.Cmd) {
	switch msg := msg.(type) {
	case auth.AuthSuccessMsg:
		m.vaultModel.SetUser(msg.User)
		m.current = VaultModel
		return m, nil
	default:
		updated, cmd := m.authModel.Update(msg)
		if authModel, ok := updated.(auth.Model); ok {
			m.authModel = authModel
		}
		return m, cmd
	}
}

func handleVaultModel(m rootModel, msg tea.Msg) (rootModel, tea.Cmd) {
	switch msg := msg.(type) {
	case vault.AddItemReqMsg:
		m.addModel.SetUser(msg.User)
		m.current = AddModel
		return m, nil
	case vault.UpdateReqMsg:
		m.updateModel.SetItem(msg.Info, msg.Content)
		m.current = UpdateModel
		return m, nil
	default:
		updated, cmd := m.vaultModel.Update(msg)
		if vaultModel, ok := updated.(vault.Model); ok {
			m.vaultModel = vaultModel
		}
		return m, cmd
	}
}

func handleAddModel(m rootModel, msg tea.Msg) (rootModel, tea.Cmd) {
	switch msg := msg.(type) {
	case add.CancelMsg:
		m.vaultModel.SetUpdateState()
		m.current = VaultModel
		return m, nil
	default:
		updated, cmd := m.addModel.Update(msg)
		if addModel, ok := updated.(add.Model); ok {
			m.addModel = addModel
		}
		return m, cmd
	}
}

func handleUpdateModel(m rootModel, msg tea.Msg) (rootModel, tea.Cmd) {
	switch msg := msg.(type) {
	case update.CancelMsg:
		m.vaultModel.SetUpdateState()
		m.current = VaultModel
		return m, nil
	default:
		updated, cmd := m.updateModel.Update(msg)
		if updateModel, ok := updated.(update.Model); ok {
			m.updateModel = updateModel
		}
		return m, cmd
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
	case UpdateModel:
		return m.updateModel.View()
	default:
		return ""
	}
}
