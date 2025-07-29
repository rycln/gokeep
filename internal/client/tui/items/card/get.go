package card

import (
	"encoding/json"
	"fmt"
	"strings"
)

func GetContentRender(content []byte) (string, error) {
	var card Card

	err := json.Unmarshal(content, &card)
	if err != nil {
		return "", err
	}

	var b strings.Builder

	b.WriteString(fmt.Sprintf("Номер: %s\n", card.CardNumber))
	b.WriteString(fmt.Sprintf("Имя держателя: %s\n", card.CardOwner))
	b.WriteString(fmt.Sprintf("Срок: %s\n", card.ExpiryDate))
	b.WriteString(fmt.Sprintf("CVV: %s\n", card.CVV))

	return b.String(), nil
}
