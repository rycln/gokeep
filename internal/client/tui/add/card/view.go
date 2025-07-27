package card

import "strings"

func (m Model) View() string {
	var b strings.Builder

	b.WriteString("Введите данные: \n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	b.WriteString("\n\nДля сохранения нажмите ENTER...\n")
	b.WriteString("\nДля отмены нажмите ESC...\n")

	return b.String()
}
