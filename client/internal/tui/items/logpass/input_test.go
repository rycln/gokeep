package logpass

import (
	"encoding/json"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/client/internal/tui/shared/messages"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialModel(t *testing.T) {
	t.Run("should initialize model with correct inputs", func(t *testing.T) {
		m := InitialModel()

		require.Len(t, m.Inputs, 4)

		assert.Equal(t, 30, m.Inputs[0].Width)
		assert.Equal(t, 32, m.Inputs[0].CharLimit)
		assert.True(t, m.Inputs[0].Focused())

		assert.Equal(t, textinput.EchoPassword, m.Inputs[3].EchoMode)
		assert.Equal(t, '•', m.Inputs[3].EchoCharacter)
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
	t.Run("successful credentials packaging", func(t *testing.T) {
		m := InitialModel()
		m.Inputs[0].SetValue("Test Credentials")
		m.Inputs[1].SetValue("Work account")
		m.Inputs[2].SetValue("test_user")
		m.Inputs[3].SetValue("secure_password123")

		cmd := m.send()
		msg := cmd()

		itemMsg, ok := msg.(messages.ItemMsg)
		require.True(t, ok)

		assert.Equal(t, "Test Credentials", itemMsg.Info.Name)
		assert.Equal(t, "Work account", itemMsg.Info.Metadata)
		assert.Equal(t, models.TypePassword, itemMsg.Info.ItemType)

		var logPass LogPass
		err := json.Unmarshal(itemMsg.Content, &logPass)
		require.NoError(t, err)
		assert.Equal(t, "test_user", logPass.Login)
		assert.Equal(t, "secure_password123", logPass.Password)

		assert.Empty(t, m.Inputs[0].Value())
	})
}

func TestModel_SetStartData(t *testing.T) {
	t.Run("successful data setup", func(t *testing.T) {
		m := InitialModel()
		testInfo := &models.ItemInfo{
			Name:     "Personal Account",
			Metadata: "Email service",
		}
		testContent := []byte(`{
			"login": "user@example.com",
			"password": "qwerty123"
		}`)

		err := m.SetStartData(testInfo, testContent)
		require.NoError(t, err)

		assert.Equal(t, "Personal Account", m.Inputs[0].Value())
		assert.Equal(t, "Email service", m.Inputs[1].Value())
		assert.Equal(t, "user@example.com", m.Inputs[2].Value())
		assert.Equal(t, "qwerty123", m.Inputs[3].Value())
	})

	t.Run("invalid json content", func(t *testing.T) {
		m := InitialModel()
		testInfo := &models.ItemInfo{
			Name:     "Test",
			Metadata: "Test",
		}

		err := m.SetStartData(testInfo, []byte("invalid json"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal error")
	})
}
