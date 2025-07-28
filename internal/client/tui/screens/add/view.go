package add

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

func (m Model) View() string {
	switch m.state {
	case SelectState:
		return m.selectView()
	case ProcessingState:
		return "Пожалуйста, подождите..."
	case AddPassword:
		return m.logpassModel.View()
	case AddCard:
		return m.cardModel.View()
	case AddText:
		return m.textModel.View()
	case ErrorState:
		return m.errorView()
	default:
		return ""
	}
}

func (m Model) selectView() string {
	s := "Выберите тип хранимой информации:\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nДля отмены нажмите ESC..."

	return s
}

func (m Model) errorView() string {
	return errorStyle.Render("Ошибка: " + m.errMsg + "\n" + "Нажмите Enter для продолжения...")
}
