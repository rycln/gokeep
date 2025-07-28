package vault

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle      = lipgloss.NewStyle().MarginLeft(2)
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
)

func (m Model) View() string {
	switch m.state {
	case ListState:
		return m.listView()
	case DetailState:
		return m.detailView()
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
	b.WriteString("Нажмите ESC для возврата к списку...")
	return b.String()
}
