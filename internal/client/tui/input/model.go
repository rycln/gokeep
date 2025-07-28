package input

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/rycln/gokeep/internal/shared/models"
)

type Form struct {
	Inputs  []textinput.Model
	Focused int
}

type (
	ItemMsg struct {
		Info    *models.ItemInfo
		Content []byte
	}
	ErrMsg    struct{ Err error }
	CancelMsg struct{}
)
