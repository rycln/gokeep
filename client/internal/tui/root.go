// Package tui implements the terminal user interface root component.
// Manages screen transitions and state between all application views.
package tui

import (
	"github.com/rycln/gokeep/client/internal/tui/screens/add"
	"github.com/rycln/gokeep/client/internal/tui/screens/auth"
	"github.com/rycln/gokeep/client/internal/tui/screens/update"
	"github.com/rycln/gokeep/client/internal/tui/screens/vault"

	tea "github.com/charmbracelet/bubbletea"
)

// model represents current active screen
type model int

// Screen type constants
const (
	AuthModel   model = iota // Authentication screen
	VaultModel               // Main vault screen
	AddModel                 // Add item screen
	UpdateModel              // Update item screen
)

// rootModel manages all application screens and transitions
type rootModel struct {
	authModel   auth.Model   // Authentication screen model
	vaultModel  vault.Model  // Main vault screen model
	addModel    add.Model    // Add item screen model
	updateModel update.Model // Update item screen model
	current     model        // Currently active screen
}

// InitialRootModel creates root model with all screen dependencies
func InitialRootModel(auth auth.Model, vault vault.Model, add add.Model, update update.Model) rootModel {
	return rootModel{
		authModel:   auth,
		vaultModel:  vault,
		addModel:    add,
		updateModel: update,
		current:     AuthModel,
	}
}

// Init initializes the root model
func (m rootModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and screen transitions
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

// handleAuthModel processes auth screen
func handleAuthModel(m rootModel, msg tea.Msg) (rootModel, tea.Cmd) {
	switch msg := msg.(type) {
	case auth.AuthSuccessMsg:
		m.vaultModel.SetUser(msg.User)
		m.current = VaultModel // Switch to vault after auth
		return m, nil
	default:
		updated, cmd := m.authModel.Update(msg)
		if authModel, ok := updated.(auth.Model); ok {
			m.authModel = authModel
		}
		return m, cmd
	}
}

// handleVaultModel processes main vault screen
func handleVaultModel(m rootModel, msg tea.Msg) (rootModel, tea.Cmd) {
	switch msg := msg.(type) {
	case vault.AddItemReqMsg:
		m.addModel.SetUser(msg.User)
		m.current = AddModel // Switch to add screen
		return m, nil
	case vault.UpdateReqMsg:
		m.updateModel.SetItem(msg.Info, msg.Content)
		m.current = UpdateModel // Switch to update screen
		return m, nil
	default:
		updated, cmd := m.vaultModel.Update(msg)
		if vaultModel, ok := updated.(vault.Model); ok {
			m.vaultModel = vaultModel
		}
		return m, cmd
	}
}

// handleAddModel processes add screen
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

// handleUpdateModel processes update screen
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

// View renders current active screen
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
