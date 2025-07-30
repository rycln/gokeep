// Package messages defines the communication protocol between TUI components.
package messages

import "github.com/rycln/gokeep/shared/models"

// Message types for TUI component communication
type (
	// ItemMsg carries item data between components
	// Contains both metadata (Info) and content
	ItemMsg struct {
		Info    *models.ItemInfo
		Content []byte
	}

	// ErrMsg transports error information between components
	ErrMsg struct{ Err error }

	// CancelMsg signals operation cancellation
	CancelMsg struct{}
)
