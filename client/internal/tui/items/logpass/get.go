package logpass

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
)

// GetContentRender formats login/password credentials for display.
// Masks sensitive information and returns formatted string.
// Returns error if content cannot be unmarshaled.
func GetContentRender(content []byte) (string, error) {
	var logPass LogPass

	err := json.Unmarshal(content, &logPass)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	var b strings.Builder

	b.WriteString(fmt.Sprintf(i18n.LogPassInputLogin+": %s\n", logPass.Login))
	b.WriteString(fmt.Sprintf(i18n.LogPassInputPassword+": %s\n", logPass.Password))

	return b.String(), nil
}
