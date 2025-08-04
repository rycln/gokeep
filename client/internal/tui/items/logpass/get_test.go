package logpass

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetContentRender(t *testing.T) {
	// Test credentials data
	testCredentials := LogPass{
		Login:    "test_user",
		Password: "secret123",
	}

	// Marshal test credentials to JSON
	validContent, err := json.Marshal(testCredentials)
	require.NoError(t, err)

	t.Run("successful credentials rendering", func(t *testing.T) {
		result, err := GetContentRender(validContent)
		require.NoError(t, err)

		expected := fmt.Sprintf(
			"%s: %s\n%s: %s\n",
			i18n.LogPassInputLogin, testCredentials.Login,
			i18n.LogPassInputPassword, testCredentials.Password,
		)

		assert.Equal(t, expected, result)
	})

	t.Run("empty content", func(t *testing.T) {
		_, err := GetContentRender([]byte{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal credentials")
	})
}

func TestLogPassStruct(t *testing.T) {
	t.Run("credentials json tags", func(t *testing.T) {
		creds := LogPass{
			Login:    "test_user",
			Password: "test_pass",
		}

		data, err := json.Marshal(creds)
		require.NoError(t, err)

		var unmarshaled LogPass
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, creds, unmarshaled)
	})
}
