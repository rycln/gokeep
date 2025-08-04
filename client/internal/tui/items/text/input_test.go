package text

import (
	"encoding/json"
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

		require.Len(t, m.Inputs, 3)

		assert.Equal(t, 30, m.Inputs[0].Width)
		assert.Equal(t, 32, m.Inputs[0].CharLimit)
		assert.True(t, m.Inputs[0].Focused())

		assert.Equal(t, 144, m.Inputs[2].Width)
		assert.Equal(t, 0, m.Inputs[2].CharLimit)
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
	t.Run("successful text data packaging", func(t *testing.T) {
		m := InitialModel()
		m.Inputs[0].SetValue("My Notes")
		m.Inputs[1].SetValue("Personal")
		m.Inputs[2].SetValue("This is a long text content")

		cmd := m.send()
		msg := cmd()

		itemMsg, ok := msg.(messages.ItemMsg)
		require.True(t, ok)

		assert.Equal(t, "My Notes", itemMsg.Info.Name)
		assert.Equal(t, "Personal", itemMsg.Info.Metadata)
		assert.Equal(t, models.TypeText, itemMsg.Info.ItemType)

		var text Text
		err := json.Unmarshal(itemMsg.Content, &text)
		require.NoError(t, err)
		assert.Equal(t, "This is a long text content", text.Content)

		assert.Empty(t, m.Inputs[0].Value())
	})
}

func TestModel_SetStartData(t *testing.T) {
	t.Run("successful data setup", func(t *testing.T) {
		m := InitialModel()
		testInfo := &models.ItemInfo{
			Name:     "Work Notes",
			Metadata: "Meeting minutes",
		}
		testContent := []byte(`{"text": "Discussion points: 1. Item one 2. Item two"}`)

		err := m.SetStartData(testInfo, testContent)
		require.NoError(t, err)

		assert.Equal(t, "Work Notes", m.Inputs[0].Value())
		assert.Equal(t, "Meeting minutes", m.Inputs[1].Value())
		assert.Equal(t, "Discussion points: 1. Item one 2. Item two", m.Inputs[2].Value())
	})

	t.Run("invalid json content", func(t *testing.T) {
		m := InitialModel()
		testInfo := &models.ItemInfo{
			Name:     "Test",
			Metadata: "Test",
		}

		err := m.SetStartData(testInfo, []byte("invalid json"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal text content")
	})
}
