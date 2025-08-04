package card

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetContentRender(t *testing.T) {
	testCard := Card{
		CardNumber: "1234 5678 9012 3456",
		CardOwner:  "John Doe",
		ExpiryDate: "12/25",
		CVV:        "123",
	}

	validContent, err := json.Marshal(testCard)
	require.NoError(t, err)

	t.Run("successful card rendering", func(t *testing.T) {
		result, err := GetContentRender(validContent)
		require.NoError(t, err)

		expected := fmt.Sprintf(
			"%s: %s\n%s: %s\n%s: %s\n%s: %s\n",
			i18n.CardInputNumber, testCard.CardNumber,
			i18n.CardInputHolderName, testCard.CardOwner,
			i18n.CardInputExpiryDate, testCard.ExpiryDate,
			i18n.CardInputCVV, testCard.CVV,
		)

		assert.Equal(t, expected, result)
	})

	t.Run("invalid JSON content", func(t *testing.T) {
		invalidContent := []byte(`{"invalid": "data"`)
		_, err := GetContentRender(invalidContent)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal card data")
	})

	t.Run("empty content", func(t *testing.T) {
		_, err := GetContentRender([]byte{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal card data")
	})
}

func TestCardStruct(t *testing.T) {
	t.Run("card json tags", func(t *testing.T) {
		card := Card{
			CardNumber: "1234 5678 9012 3456",
			CardOwner:  "John Doe",
			ExpiryDate: "12/25",
			CVV:        "123",
		}

		data, err := json.Marshal(card)
		require.NoError(t, err)

		var unmarshaled Card
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, card, unmarshaled)
	})
}
