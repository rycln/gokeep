package update

import (
	"fmt"

	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/rycln/gokeep/client/internal/tui/shared/styles"
)

// View renders the current update screen based on state.
// Returns formatted UI with appropriate styling and localization.
func (m Model) View() string {
	switch m.state {
	case LoadState:
		return i18n.CommonPressAnyKey
	case ProcessingState:
		return i18n.CommonWait
	case UpdatePassword:
		return m.logpassModel.View()
	case UpdateCard:
		return m.cardModel.View()
	case UpdateText:
		return m.textModel.View()
	case UpdateBinary:
		return m.binModel.View()
	case ErrorState:
		return m.errorView()
	default:
		return ""
	}
}

// errorView renders error messages with consistent styling.
// Formats error message using shared error style and localization.
func (m Model) errorView() string {
	return styles.ErrorStyle.Render(
		fmt.Sprintf(i18n.CommonError, m.errMsg),
	)
}
