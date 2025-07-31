package input

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/rycln/gokeep/client/internal/tui/shared/messages"
	"github.com/rycln/gokeep/client/internal/tui/shared/styles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForm_Init(t *testing.T) {
	t.Run("should return blink command", func(t *testing.T) {
		form := Form{
			Inputs: make([]textinput.Model, 1),
		}
		cmd := form.Init()
		assert.NotNil(t, cmd)
	})
}

func TestForm_Update(t *testing.T) {
	t.Run("should handle Esc to cancel", func(t *testing.T) {
		form := Form{
			Inputs: make([]textinput.Model, 1),
		}
		_, cmd := form.Update(tea.KeyMsg{Type: tea.KeyEsc})
		require.NotNil(t, cmd)

		msg := cmd()
		_, ok := msg.(messages.CancelMsg)
		assert.True(t, ok)
	})

	t.Run("should move focus down", func(t *testing.T) {
		form := Form{
			Inputs:  make([]textinput.Model, 3),
			Focused: 0,
		}

		for i := range form.Inputs {
			t := textinput.New()
			form.Inputs[i] = t
		}

		form.Inputs[0].Focus()

		updated, cmd := form.Update(tea.KeyMsg{Type: tea.KeyDown})
		if formModel, ok := updated.(Form); ok {
			form = formModel
		}
		assert.NotNil(t, cmd)
		assert.Equal(t, 1, form.Focused)
	})

	t.Run("should wrap focus from bottom to top", func(t *testing.T) {
		form := Form{
			Inputs:  make([]textinput.Model, 2),
			Focused: 1,
		}

		for i := range form.Inputs {
			t := textinput.New()
			form.Inputs[i] = t
		}

		updated, cmd := form.Update(tea.KeyMsg{Type: tea.KeyDown})
		if formModel, ok := updated.(Form); ok {
			form = formModel
		}

		assert.NotNil(t, cmd)
		assert.Equal(t, 0, form.Focused)
	})

	t.Run("should move focus up", func(t *testing.T) {
		form := Form{
			Inputs:  make([]textinput.Model, 2),
			Focused: 1,
		}

		for i := range form.Inputs {
			t := textinput.New()
			form.Inputs[i] = t
		}

		updated, cmd := form.Update(tea.KeyMsg{Type: tea.KeyUp})
		if formModel, ok := updated.(Form); ok {
			form = formModel
		}

		assert.NotNil(t, cmd)
		assert.Equal(t, 0, form.Focused)
	})

	t.Run("should wrap focus from top to bottom", func(t *testing.T) {
		form := Form{
			Inputs:  make([]textinput.Model, 2),
			Focused: 0,
		}

		for i := range form.Inputs {
			t := textinput.New()
			form.Inputs[i] = t
		}

		updated, cmd := form.Update(tea.KeyMsg{Type: tea.KeyUp})
		if formModel, ok := updated.(Form); ok {
			form = formModel
		}
		assert.NotNil(t, cmd)
		assert.Equal(t, 1, form.Focused)
	})

	t.Run("should delegate input to focused field", func(t *testing.T) {
		form := Form{
			Inputs:  make([]textinput.Model, 1),
			Focused: 0,
		}
		for i := range form.Inputs {
			t := textinput.New()
			form.Inputs[i] = t
		}

		form.Inputs[0].Focus()

		updated, cmd := form.Update(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'a'},
		})
		if formModel, ok := updated.(Form); ok {
			form = formModel
		}

		assert.NotNil(t, cmd)
		assert.Equal(t, "a", form.Inputs[0].Value())
	})
}

func TestForm_updateFocus(t *testing.T) {
	t.Run("should update focus styles", func(t *testing.T) {
		form := Form{
			Inputs:  make([]textinput.Model, 2),
			Focused: 0,
		}

		for i := range form.Inputs {
			t := textinput.New()
			form.Inputs[i] = t
		}

		cmd := form.updateFocus()
		assert.NotNil(t, cmd)

		assert.True(t, form.Inputs[0].Focused())
		assert.Equal(t, styles.FocusedStyle, form.Inputs[0].PromptStyle)
		assert.Equal(t, styles.FocusedStyle, form.Inputs[0].TextStyle)

		assert.False(t, form.Inputs[1].Focused())
		assert.Equal(t, styles.NoStyle, form.Inputs[1].PromptStyle)
		assert.Equal(t, styles.NoStyle, form.Inputs[1].TextStyle)
	})
}

func TestForm_View(t *testing.T) {
	t.Run("should render all inputs and instructions", func(t *testing.T) {
		form := Form{
			Inputs: make([]textinput.Model, 2),
		}

		for i := range form.Inputs {
			t := textinput.New()
			t.Width = 10

			switch i {
			case 0:
				t.Placeholder = "First"
			case 1:
				t.Placeholder = "Second"
			}

			form.Inputs[i] = t
		}

		view := form.View()
		assert.Contains(t, view, i18n.InputDataPrompt)
		assert.Contains(t, view, "First")
		assert.Contains(t, view, "Second")
		assert.Contains(t, view, i18n.CommonPressEnter)
		assert.Contains(t, view, i18n.CommonPressESC)
	})
}
