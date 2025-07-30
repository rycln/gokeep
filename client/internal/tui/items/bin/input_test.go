package bin

import (
	"encoding/json"
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/client/internal/tui/shared/messages"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialModel(t *testing.T) {
	t.Run("should initialize model with correct inputs", func(t *testing.T) {
		m := InitialModel()

		assert.Len(t, m.Inputs, 3)
		assert.Equal(t, 30, m.Inputs[0].Width)
		assert.Equal(t, 32, m.Inputs[0].CharLimit)
		assert.True(t, m.Inputs[0].Focused())
		assert.Equal(t, 100, m.Inputs[2].CharLimit)
	})
}

func TestModel_Update(t *testing.T) {
	t.Run("enter key should trigger send command", func(t *testing.T) {
		m := InitialModel()
		_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		assert.NotNil(t, cmd)
	})
}

func TestModel_send(t *testing.T) {
	t.Run("successful file read and processing", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "testfile")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		testContent := []byte("test content")
		err = os.WriteFile(tmpFile.Name(), testContent, 0644)
		require.NoError(t, err)

		m := InitialModel()
		m.Inputs[0].SetValue("test name")
		m.Inputs[1].SetValue("test metadata")
		m.Inputs[2].SetValue(tmpFile.Name())

		cmd := m.send()
		msg := cmd()

		itemMsg, ok := msg.(messages.ItemMsg)
		require.True(t, ok)

		assert.Equal(t, "test name", itemMsg.Info.Name)
		assert.Equal(t, "test metadata", itemMsg.Info.Metadata)
		assert.Equal(t, models.TypeBinary, itemMsg.Info.ItemType)

		var binFile BinFile
		err = json.Unmarshal(itemMsg.Content, &binFile)
		require.NoError(t, err)
		assert.Equal(t, testContent, binFile.Data)
	})

	t.Run("file read error", func(t *testing.T) {
		m := InitialModel()
		m.Inputs[2].SetValue("/nonexistent/file")

		cmd := m.send()
		msg := cmd()

		errMsg, ok := msg.(messages.ErrMsg)
		require.True(t, ok)
		assert.Error(t, errMsg.Err)
		assert.Contains(t, errMsg.Err.Error(), "file read error")
	})
}

func TestModel_SetStartData(t *testing.T) {
	t.Run("successful data setup", func(t *testing.T) {
		m := InitialModel()
		testInfo := &models.ItemInfo{
			Name:     "test name",
			Metadata: "test metadata",
		}
		testContent := []byte(`{"data":"dGVzdCBjb250ZW50"}`)

		err := m.SetStartData(testInfo, testContent)
		require.NoError(t, err)

		assert.Equal(t, "test name", m.Inputs[0].Value())
		assert.Equal(t, "test metadata", m.Inputs[1].Value())
	})

	t.Run("invalid json content", func(t *testing.T) {
		m := InitialModel()
		testInfo := &models.ItemInfo{
			Name:     "test name",
			Metadata: "test metadata",
		}

		err := m.SetStartData(testInfo, []byte("invalid json"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "json unmarshal error")
	})
}
