package update

import (
	"fmt"

	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/rycln/gokeep/client/internal/tui/shared/styles"
)

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

func (m Model) errorView() string {
	return styles.ErrorStyle.Render(fmt.Sprintf(i18n.CommonError, m.errMsg))
}
