package card

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

		require.Len(t, m.Inputs, 6)

		assert.Equal(t, 40, m.Inputs[0].Width)
		assert.Equal(t, 32, m.Inputs[0].CharLimit)
		assert.True(t, m.Inputs[0].Focused())

		assert.Equal(t, 16, m.Inputs[2].CharLimit)
		assert.Equal(t, 40, m.Inputs[3].CharLimit)
		assert.Equal(t, 5, m.Inputs[4].CharLimit)
		assert.Equal(t, textinput.EchoPassword, m.Inputs[5].EchoMode)
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
	t.Run("successful card data packaging", func(t *testing.T) {
		m := InitialModel()
		m.Inputs[0].SetValue("Test Card")
		m.Inputs[1].SetValue("Personal")
		m.Inputs[2].SetValue("4111111111111111")
		m.Inputs[3].SetValue("John Doe")
		m.Inputs[4].SetValue("12/25")
		m.Inputs[5].SetValue("123")

		cmd := m.send()
		msg := cmd()

		itemMsg, ok := msg.(messages.ItemMsg)
		require.True(t, ok)

		assert.Equal(t, "Test Card", itemMsg.Info.Name)
		assert.Equal(t, "Personal", itemMsg.Info.Metadata)
		assert.Equal(t, models.TypeCard, itemMsg.Info.ItemType)

		var card Card
		err := json.Unmarshal(itemMsg.Content, &card)
		require.NoError(t, err)
		assert.Equal(t, "4111111111111111", card.CardNumber)
		assert.Equal(t, "John Doe", card.CardOwner)
		assert.Equal(t, "12/25", card.ExpiryDate)
		assert.Equal(t, "123", card.CVV)

		assert.Empty(t, m.Inputs[0].Value())
	})
}

func TestModel_SetStartData(t *testing.T) {
	t.Run("successful data setup", func(t *testing.T) {
		m := InitialModel()
		testInfo := &models.ItemInfo{
			Name:     "Test Card",
			Metadata: "Work",
		}
		testContent := []byte(`{
			"card_number": "5555555555554444",
			"card_owner": "Jane Smith",
			"expiry_date": "06/24",
			"cvv": "456"
		}`)

		err := m.SetStartData(testInfo, testContent)
		require.NoError(t, err)

		assert.Equal(t, "Test Card", m.Inputs[0].Value())
		assert.Equal(t, "Work", m.Inputs[1].Value())
		assert.Equal(t, "5555555555554444", m.Inputs[2].Value())
		assert.Equal(t, "Jane Smith", m.Inputs[3].Value())
		assert.Equal(t, "06/24", m.Inputs[4].Value())
		assert.Equal(t, "456", m.Inputs[5].Value())
	})

	t.Run("invalid json content", func(t *testing.T) {
		m := InitialModel()
		testInfo := &models.ItemInfo{
			Name:     "Test Card",
			Metadata: "Work",
		}

		err := m.SetStartData(testInfo, []byte("invalid json"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal card data")
	})
}
