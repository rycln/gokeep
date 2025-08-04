package bin

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
)

// UploadFile saves binary content to the specified path
// Returns success message or error if operation fails
func UploadFile(path string, content []byte) (string, error) {
	var binary BinFile

	err := json.Unmarshal(content, &binary)
	if err != nil {
		return "", fmt.Errorf("json unmarshal failed: %w", err)
	}

	err = os.WriteFile(path, binary.Data, 0644)
	if err != nil {
		return "", fmt.Errorf("file write failed: %w", err)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf(i18n.BinInputSuccess+"%s\n", path))

	return b.String(), nil
}
