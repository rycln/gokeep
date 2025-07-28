package input

import (
	"strings"
)

func (m Form) View() string {
	var b strings.Builder

	b.WriteString("Введите данные: \n\n")

	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	b.WriteString("\n\nДля сохранения нажмите ENTER...\n")
	b.WriteString("\nДля отмены нажмите ESC...\n")

	return b.String()
}
