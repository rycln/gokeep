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
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

func (m Model) View() string {
	switch m.mode {
	case "list":
		return m.listView()
	case "detail":
		return m.detailView()
	default:
		return ""
	}
}

func (m Model) listView() string {
	return m.list.View() + "\nPress 'n' to add new item, 'q' to quit"
}

func (m Model) detailView() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Selected: %s\n", m.selected.Name))
	b.WriteString(fmt.Sprintf("Type: %s\n", m.selected.ItemType))
	b.WriteString(fmt.Sprintf("Content:\n%s\n\n", m.content))
	b.WriteString("Press ESC to return to list")
	return b.String()
}
