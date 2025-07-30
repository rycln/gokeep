// Package update implements the item modification interface.
// Handles loading and updating of different item types.
package update

import (
	"context"
	"time"

	"github.com/rycln/gokeep/client/internal/tui/items/bin"
	"github.com/rycln/gokeep/client/internal/tui/items/card"
	"github.com/rycln/gokeep/client/internal/tui/items/logpass"
	"github.com/rycln/gokeep/client/internal/tui/items/text"
	"github.com/rycln/gokeep/shared/models"
)

// state represents the current update screen state
type state int

// Update item screen states
const (
	LoadState       state = iota // Initial loading state
	UpdatePassword               // Updating login credentials
	UpdateCard                   // Updating payment card
	UpdateText                   // Updating text content
	UpdateBinary                 // Updating binary file
	ProcessingState              // Processing update operation
	ErrorState                   // Error display state
)

// itemUpdater defines interface for updating items
type itemUpdater interface {
	Update(context.Context, *models.ItemInfo, []byte) error
}

// Message types for update operations
type (
	// InitSuccessMsg signals successful initialization
	InitSuccessMsg struct{ state state }
	// UpdateSuccessMsg indicates successful item update
	UpdateSuccessMsg struct{}
	// ErrorMsg contains operation error details
	ErrorMsg struct{ Err error }
	// CancelMsg signals operation cancellation
	CancelMsg struct{}
)

// Model manages the update item screen state
type Model struct {
	state        state            // Current screen state
	errMsg       string           // Last error message
	logpassModel logpass.Model    // Login/password form
	cardModel    card.Model       // Card form
	textModel    text.Model       // Text form
	binModel     bin.Model        // Binary form
	info         *models.ItemInfo // Item metadata
	content      []byte           // Current item content
	service      itemUpdater      // Item storage service
	timeout      time.Duration    // Operation timeout
}

// InitialModel creates new update item model with dependencies
func InitialModel(service itemUpdater, timeout time.Duration) Model {
	return Model{
		state:        LoadState,
		logpassModel: logpass.InitialModel(),
		cardModel:    card.InitialModel(),
		textModel:    text.InitialModel(),
		binModel:     bin.InitialModel(),
		service:      service,
		timeout:      timeout,
	}
}

// SetItem prepares the model for updating a specific item
func (m *Model) SetItem(info *models.ItemInfo, content []byte) {
	m.state = LoadState
	m.info = info
	m.content = content
}
