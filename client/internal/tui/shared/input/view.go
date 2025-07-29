package input

import (
	"strings"

	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
)

func (m Form) View() string {
	var b strings.Builder
	b.WriteString(i18n.InputDataPrompt)
	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}
	b.WriteString("\n\n" + i18n.CommonPressEnter + "\n")
	b.WriteString("\n" + i18n.CommonPressESC + "\n")
	return b.String()
}
