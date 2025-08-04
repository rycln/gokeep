package vault

import (
	"fmt"
	"strings"

	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/rycln/gokeep/client/internal/tui/shared/styles"
)

// View renders the current state of the vault interface.
// Returns different views based on the current state machine state.
func (m Model) View() string {
	switch m.state {
	case UpdateState:
		return i18n.CommonPressAnyKey
	case ProcessingState:
		return i18n.CommonWait
	case ListState:
		return m.listView()
	case DetailState:
		return m.detailView()
	case BinaryInputState:
		return fmt.Sprintf(i18n.InputSavePathPrompt, m.input)
	case ErrorState:
		return styles.ErrorStyle.Render(fmt.Sprintf(i18n.CommonError, m.errMsg))
	default:
		return ""
	}
}

// detailView renders the detailed view of a selected item.
// Shows all item metadata and formatted content with action prompts.
func (m Model) detailView() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(i18n.VaultObjectTitle, m.selected.Name) + "\n")
	b.WriteString(fmt.Sprintf(i18n.VaultTypeTitle, m.selected.ItemType) + "\n")
	b.WriteString(fmt.Sprintf(i18n.VaultDescTitle, m.selected.Metadata) + "\n")
	b.WriteString(fmt.Sprintf(i18n.VaultUpdatedTitle, m.selected.UpdatedAt.String()))
	if m.selected.Content != "" {
		b.WriteString(m.selected.Content + "\n")
	}
	b.WriteString(i18n.VaultActions)
	return b.String()
}

// listView renders the main items list view.
// Uses the bubbletea list component for consistent list rendering.
func (m Model) listView() string {
	return m.list.View()
}
