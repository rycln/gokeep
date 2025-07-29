package vault

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	titleStyle      = lipgloss.NewStyle().MarginLeft(2)
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
)

func (m Model) View() string {
	switch m.state {
	case StartState:
		return "Нажмите любую клавишу..."
	case ProcessingState:
		return "Пожалуйста, подождите..."
	case ListState:
		return m.listView()
	case DetailState:
		return m.detailView()
	case ErrorState:
		return errorStyle.Render("Ошибка: " + m.errMsg + "\n" + "Нажмите Enter для продолжения...")
	default:
		return ""
	}
}

func (m Model) listView() string {
	return m.list.View()
}

func (m Model) detailView() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Объект: %s\n", m.selected.Name))
	b.WriteString(fmt.Sprintf("Тип: %s\n", m.selected.ItemType))
	if m.selected.Content != "" {
		b.WriteString(m.selected.Content)
	}
	b.WriteString("Нажмите ENTER для загрузки данных...\n\n")
	b.WriteString("Нажмите ESC для возврата к списку...")
	return b.String()
}
