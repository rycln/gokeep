package input

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type Form struct {
	Inputs  []textinput.Model
	Focused int
}
