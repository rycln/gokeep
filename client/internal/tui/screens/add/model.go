// Package add implements the item addition interface.
// Handles selection and creation of different item types.
package add

import (
	"context"
	"time"

	"github.com/rycln/gokeep/client/internal/tui/items/bin"
	"github.com/rycln/gokeep/client/internal/tui/items/card"
	"github.com/rycln/gokeep/client/internal/tui/items/logpass"
	"github.com/rycln/gokeep/client/internal/tui/items/text"
	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/rycln/gokeep/shared/models"
)

// state represents the current screen state
type state int

// Add item screen states
const (
	SelectState     state = iota // Item type selection
	AddPassword                  // Adding login credentials
	AddCard                      // Adding payment card
	AddText                      // Adding text content
	AddBinary                    // Adding binary file
	ProcessingState              // Processing operation
	ErrorState                   // Error display
)

// choice represents available item types
type choice string

// itemAdder defines interface for adding items
type itemAdder interface {
	Add(context.Context, *models.ItemInfo, []byte) error
}

// Message types for add operations
type (
	// AddSuccessMsg indicates successful item addition
	AddSuccessMsg struct{}
	// ErrorMsg contains operation error details
	ErrorMsg struct{ Err error }
	// CancelMsg signals operation cancellation
	CancelMsg struct{}
)

// Model manages the add item screen state
type Model struct {
	state        state         // Current screen state
	choices      []choice      // Available item types
	cursor       int           // Selection cursor position
	errMsg       string        // Last error message
	logpassModel logpass.Model // Login/password form
	cardModel    card.Model    // Card form
	textModel    text.Model    // Text form
	binModel     bin.Model     // Binary form
	user         *models.User  // Current user
	service      itemAdder     // Item storage service
	timeout      time.Duration // Operation timeout
}

// InitialModel creates new add item model with dependencies
func InitialModel(service itemAdder, timeout time.Duration) Model {
	return Model{
		state:        SelectState,
		choices:      []choice{i18n.AddPassword, i18n.AddCard, i18n.AddText, i18n.AddBinary},
		logpassModel: logpass.InitialModel(),
		cardModel:    card.InitialModel(),
		textModel:    text.InitialModel(),
		binModel:     bin.InitialModel(),
		service:      service,
		timeout:      timeout,
	}
}

// SetUser updates the current authenticated user
func (m *Model) SetUser(user *models.User) {
	m.user = user
}
