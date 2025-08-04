package text

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetContentRender(t *testing.T) {
	t.Run("successful text rendering", func(t *testing.T) {
		testText := Text{
			Content: "This is a sample text content\nwith multiple lines",
		}

		content, err := json.Marshal(testText)
		require.NoError(t, err)

		result, err := GetContentRender(content)
		require.NoError(t, err)

		expected := fmt.Sprintf("%s:\n%s\n", i18n.TextInputContent, testText.Content)
		assert.Equal(t, expected, result)
	})

	t.Run("empty content", func(t *testing.T) {
		_, err := GetContentRender([]byte{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal text content")
	})

	t.Run("empty text content", func(t *testing.T) {
		testText := Text{
			Content: "",
		}

		content, err := json.Marshal(testText)
		require.NoError(t, err)

		result, err := GetContentRender(content)
		require.NoError(t, err)

		expected := fmt.Sprintf("%s:\n%s\n", i18n.TextInputContent, "")
		assert.Equal(t, expected, result)
	})
}

func TestTextStruct(t *testing.T) {
	t.Run("text json tags", func(t *testing.T) {
		testText := Text{
			Content: "Test content",
		}

		data, err := json.Marshal(testText)
		require.NoError(t, err)

		var unmarshaled Text
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, testText, unmarshaled)
	})
}
