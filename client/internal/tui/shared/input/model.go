// Package input provides form input components and management.
package input

import (
	"github.com/charmbracelet/bubbles/textinput"
)

// Form manages a collection of text input fields.
// Tracks focus state and provides unified input handling.
type Form struct {
	Inputs  []textinput.Model // Collection of text input fields
	Focused int               // Index of currently focused field
}
