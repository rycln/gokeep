package logpass

import (
	"encoding/json"
	"fmt"
	"strings"
)

func GetRender(content []byte) (string, error) {
	var logPass LogPass

	err := json.Unmarshal(content, &logPass)
	if err != nil {
		return "", err
	}

	var b strings.Builder

	b.WriteString(fmt.Sprintf("Логин: %s\n", logPass.Login))
	b.WriteString(fmt.Sprintf("Пароль: %s\n", logPass.Password))

	return b.String(), nil
}
