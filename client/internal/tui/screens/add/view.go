package add

import (
	"fmt"
	"strings"

	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/rycln/gokeep/client/internal/tui/shared/styles"
)

func (m Model) View() string {
	switch m.state {
	case SelectState:
		return m.selectView()
	case ProcessingState:
		return i18n.CommonWait
	case AddPassword:
		return m.logpassModel.View()
	case AddCard:
		return m.cardModel.View()
	case AddText:
		return m.textModel.View()
	case AddBinary:
		return m.binModel.View()
	case ErrorState:
		return m.errorView()
	default:
		return ""
	}
}

func (m Model) selectView() string {
	var b strings.Builder
	b.WriteString(i18n.AddSelectPrompt)
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf(i18n.AddChoiceTemplate, cursor, choice))
	}
	b.WriteString("\n" + i18n.CommonPressESC)
	return b.String()
}

func (m Model) errorView() string {
	return styles.ErrorStyle.Render(fmt.Sprintf(i18n.CommonError, m.errMsg))
}
