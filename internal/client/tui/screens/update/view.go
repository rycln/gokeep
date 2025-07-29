package update

import (
	"github.com/charmbracelet/lipgloss"
)

var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

func (m Model) View() string {
	switch m.state {
	case LoadState:
		return "Нажмите любую клавишу..."
	case ProcessingState:
		return "Пожалуйста, подождите..."
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
	return errorStyle.Render("Ошибка: " + m.errMsg + "\n" + "Нажмите Enter для продолжения...")
}
