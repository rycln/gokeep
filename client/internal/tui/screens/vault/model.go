// Package vault implements the main vault screen for item management.
// Handles listing, viewing and managing all stored secret items.
package vault

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/rycln/gokeep/client/internal/tui/shared/styles"
	"github.com/rycln/gokeep/shared/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// state represents current view state of vault screen
type state int

// Screen state constants
const (
	ListState        state = iota // Items list view
	DetailState                   // Item detail view
	UpdateState                   // Updating items state
	BinaryInputState              // Binary input processing
	ProcessingState               // Background operation in progress
	ErrorState                    // Error display state
)

// Message types for vault screen communication
type (
	// AddItemReqMsg requests showing add item screen
	AddItemReqMsg struct{ User *models.User }

	// UpdateReqMsg requests showing update screen for item
	UpdateReqMsg struct {
		Info    *models.ItemInfo
		Content []byte
	}

	// DeleteSuccessMsg confirms successful item deletion
	DeleteSuccessMsg struct{}

	// ItemsMsg delivers list of items for display
	ItemsMsg struct{ Items []itemRender }

	// ContentMsg delivers item content for detail view
	ContentMsg struct{ Content string }

	// ErrorMsg delivers error information
	ErrorMsg struct{ Err error }
)

// itemGetter defines interface for reading items
type itemGetter interface {
	List(context.Context, models.UserID) ([]models.ItemInfo, error)
	GetContent(context.Context, models.ItemID) ([]byte, error)
}

// itemDeleter defines interface for deleting items
type itemDeleter interface {
	Delete(context.Context, models.ItemID) error
}

// itemService combines item management interfaces
type itemService interface {
	itemGetter
	itemDeleter
}

// itemRender represents formatted item for display
type itemRender struct {
	ID        models.ItemID   // Unique item identifier
	ItemType  models.ItemType // Type of item (text, card etc)
	Name      string          // Display name
	Metadata  string          // Additional description
	UpdatedAt time.Time       // Last modification time
	Content   string          // Formatted content
}

// FilterValue implements list.Item interface for filtering
func (i itemRender) FilterValue() string { return i.Name }

// Title implements list.Item interface for display
func (i itemRender) Title() string { return i.Name }

// Description implements list.Item interface for display
func (i itemRender) Description() string {
	return fmt.Sprintf(i18n.VaultTypeTitle+"\n"+i18n.VaultDescTitle, i.ItemType, i.Metadata)
}

// Model represents vault screen state and components
type Model struct {
	state    state         // Current view state
	input    string        // User input buffer
	selected *itemRender   // Currently selected item
	items    []itemRender  // List of all items
	list     list.Model    // List UI component
	errMsg   string        // Last error message
	service  itemService   // Item service interface
	user     *models.User  // Current authenticated user
	timeout  time.Duration // UI message timeout
}

// InitialModel creates new vault model with dependencies
func InitialModel(service itemService, timeout time.Duration) Model {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("62")).
		Foreground(lipgloss.Color("229"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("245"))
	delegate.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("n"),
				key.WithHelp("n", i18n.VaultAddItemHelp),
			),
			key.NewBinding(
				key.WithKeys("u"),
				key.WithHelp("u", i18n.VaultUpdateHelp),
			),
		}
	}

	l := list.New([]list.Item{}, delegate, 20, 20)
	l.Title = i18n.VaultTitle
	l.Styles.Title = styles.TitleStyle
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.StatusMessageLifetime = timeout
	l.SetStatusBarItemName(i18n.VaultListTitleNameSingular, i18n.VaultListTitleNamePlural)

	return Model{
		state:   UpdateState,
		list:    l,
		service: service,
		timeout: timeout,
	}
}

// SetUser updates current authenticated user
func (m *Model) SetUser(user *models.User) {
	m.user = user
}

// SetUpdateState resets view to update items list
func (m *Model) SetUpdateState() {
	m.state = UpdateState
}
