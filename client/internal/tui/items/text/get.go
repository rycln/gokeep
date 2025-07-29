package text

import (
	"encoding/json"
	"fmt"
	"strings"
)

func GetContentRender(content []byte) (string, error) {
	var text Text

	err := json.Unmarshal(content, &text)
	if err != nil {
		return "", err
	}

	var b strings.Builder

	b.WriteString(fmt.Sprintf("Текст: %s\n", text.Content))

	return b.String(), nil
}
