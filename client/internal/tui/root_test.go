package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/client/internal/tui/screens/add"
	"github.com/rycln/gokeep/client/internal/tui/screens/auth"
	"github.com/rycln/gokeep/client/internal/tui/screens/update"
	"github.com/rycln/gokeep/client/internal/tui/screens/vault"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialRootModel(t *testing.T) {
	t.Run("should initialize with auth screen", func(t *testing.T) {
		authModel := auth.Model{}
		vaultModel := vault.Model{}
		addModel := add.Model{}
		updateModel := update.Model{}

		model := InitialRootModel(authModel, vaultModel, addModel, updateModel)

		assert.Equal(t, AuthModel, model.current)
		assert.Equal(t, authModel, model.authModel)
		assert.Equal(t, vaultModel, model.vaultModel)
		assert.Equal(t, addModel, model.addModel)
		assert.Equal(t, updateModel, model.updateModel)
	})
}

func TestRootModel_Update(t *testing.T) {
	t.Run("should transition from auth to vault on success", func(t *testing.T) {
		user := &models.User{ID: "user123"}
		authModel := auth.Model{}
		vaultModel := vault.Model{}
		model := InitialRootModel(authModel, vaultModel, add.Model{}, update.Model{})

		updated, cmd := model.Update(auth.AuthSuccessMsg{User: user})
		require.Nil(t, cmd)

		rootModel, ok := updated.(rootModel)
		require.True(t, ok)
		assert.Equal(t, VaultModel, rootModel.current)
	})

	t.Run("should transition from vault to add on add request", func(t *testing.T) {
		user := &models.User{ID: "user123"}
		vaultModel := vault.Model{}
		model := InitialRootModel(auth.Model{}, vaultModel, add.Model{}, update.Model{})
		model.current = VaultModel

		updated, cmd := model.Update(vault.AddItemReqMsg{User: user})
		require.Nil(t, cmd)

		rootModel, ok := updated.(rootModel)
		require.True(t, ok)
		assert.Equal(t, AddModel, rootModel.current)
	})

	t.Run("should transition from vault to update on update request", func(t *testing.T) {
		vaultModel := vault.Model{}
		model := InitialRootModel(auth.Model{}, vaultModel, add.Model{}, update.Model{})
		model.current = VaultModel

		itemInfo := &models.ItemInfo{ID: "item123"}
		content := []byte("content")
		updated, cmd := model.Update(vault.UpdateReqMsg{Info: itemInfo, Content: content})
		require.Nil(t, cmd)

		rootModel, ok := updated.(rootModel)
		require.True(t, ok)
		assert.Equal(t, UpdateModel, rootModel.current)
	})

	t.Run("should return to vault from add on cancel", func(t *testing.T) {
		vaultModel := vault.Model{}
		model := InitialRootModel(auth.Model{}, vaultModel, add.Model{}, update.Model{})
		model.current = AddModel

		updated, cmd := model.Update(add.CancelMsg{})
		require.Nil(t, cmd)

		rootModel, ok := updated.(rootModel)
		require.True(t, ok)
		assert.Equal(t, VaultModel, rootModel.current)
	})

	t.Run("should return to vault from update on cancel", func(t *testing.T) {
		vaultModel := vault.Model{}
		model := InitialRootModel(auth.Model{}, vaultModel, add.Model{}, update.Model{})
		model.current = UpdateModel

		updated, cmd := model.Update(update.CancelMsg{})
		require.Nil(t, cmd)

		rootModel, ok := updated.(rootModel)
		require.True(t, ok)
		assert.Equal(t, VaultModel, rootModel.current)
	})

	t.Run("should delegate update to current screen", func(t *testing.T) {
		authModel := auth.Model{}
		model := InitialRootModel(authModel, vault.Model{}, add.Model{}, update.Model{})

		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		assert.NotNil(t, cmd)
	})
}

func TestRootModel_View(t *testing.T) {
	t.Run("should render auth screen when active", func(t *testing.T) {
		authModel := auth.Model{}
		model := InitialRootModel(authModel, vault.Model{}, add.Model{}, update.Model{})
		model.current = AuthModel

		view := model.View()
		assert.Equal(t, authModel.View(), view)
	})

	t.Run("should render vault screen when active", func(t *testing.T) {
		vaultModel := vault.Model{}
		model := InitialRootModel(auth.Model{}, vaultModel, add.Model{}, update.Model{})
		model.current = VaultModel

		view := model.View()
		assert.Equal(t, vaultModel.View(), view)
	})

	t.Run("should render add screen when active", func(t *testing.T) {
		addModel := add.Model{}
		model := InitialRootModel(auth.Model{}, vault.Model{}, addModel, update.Model{})
		model.current = AddModel

		view := model.View()
		assert.Equal(t, addModel.View(), view)
	})

	t.Run("should render update screen when active", func(t *testing.T) {
		updateModel := update.Model{}
		model := InitialRootModel(auth.Model{}, vault.Model{}, add.Model{}, updateModel)
		model.current = UpdateModel

		view := model.View()
		assert.Equal(t, updateModel.View(), view)
	})
}
