package text

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
)

// GetContentRender formats and displays text content from stored JSON data.
func GetContentRender(content []byte) (string, error) {
	var text Text

	err := json.Unmarshal(content, &text)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal text content: %w", err)
	}

	var b strings.Builder

	b.WriteString(fmt.Sprintf(i18n.TextInputContent+":\n%s\n", text.Content))

	return b.String(), nil
}
