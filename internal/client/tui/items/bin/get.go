package bin

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func UploadFile(path string, content []byte) (string, error) {
	var binary BinFile

	err := json.Unmarshal(content, &binary)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(path, binary.Data, 0644)
	if err != nil {
		return "", err
	}

	var b strings.Builder

	b.WriteString(fmt.Sprintf("Файл сохранен успешно по пути: %s\n", path))

	return b.String(), nil
}
