package card

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
)

// GetContentRender formats card details for display from JSON content
// Returns formatted string with card information or error if unmarshaling fails
func GetContentRender(content []byte) (string, error) {
	var card Card

	err := json.Unmarshal(content, &card)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal card data: %w", err)
	}

	var b strings.Builder

	b.WriteString(fmt.Sprintf(i18n.CardInputNumber+": %s\n", card.CardNumber))
	b.WriteString(fmt.Sprintf(i18n.CardInputHolderName+": %s\n", card.CardOwner))
	b.WriteString(fmt.Sprintf(i18n.CardInputExpiryDate+": %s\n", card.ExpiryDate))
	b.WriteString(fmt.Sprintf(i18n.CardInputCVV+": %s\n", card.CVV))

	return b.String(), nil
}
