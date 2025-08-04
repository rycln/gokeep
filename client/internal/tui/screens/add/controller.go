package add

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/client/internal/tui/items/bin"
	"github.com/rycln/gokeep/client/internal/tui/items/card"
	"github.com/rycln/gokeep/client/internal/tui/items/logpass"
	"github.com/rycln/gokeep/client/internal/tui/items/text"
	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/rycln/gokeep/client/internal/tui/shared/messages"
	"github.com/rycln/gokeep/shared/models"
)

// Init initializes the add item model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles all messages and state transitions for item creation
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.ItemMsg:
		info := msg.Info
		info.UserID = m.user.ID
		m.state = ProcessingState
		return m, m.add(info, msg.Content)
	case messages.ErrMsg:
		m.errMsg = msg.Err.Error()
		m.state = ErrorState
		return m, nil
	case messages.CancelMsg:
		m.state = SelectState
		return m, nil
	default:
		switch m.state {
		case AddPassword:
			return handleItemModelUpdate(m, msg, &m.logpassModel)
		case AddCard:
			return handleItemModelUpdate(m, msg, &m.cardModel)
		case AddText:
			return handleItemModelUpdate(m, msg, &m.textModel)
		case AddBinary:
			return handleItemModelUpdate(m, msg, &m.binModel)
		default:
			switch m.state {
			case SelectState:
				return handleSelectState(m, msg)
			case ProcessingState:
				return handleProcessingState(m, msg)
			case ErrorState:
				return handleErrorState(m, msg)
			default:
				return m, nil
			}
		}
	}
}

// handleItemModelUpdate delegates updates to specific item model
func handleItemModelUpdate(m Model, msg tea.Msg, model tea.Model) (Model, tea.Cmd) {
	updated, cmd := model.Update(msg)
	switch v := updated.(type) {
	case logpass.Model:
		m.logpassModel = v
	case card.Model:
		m.cardModel = v
	case text.Model:
		m.textModel = v
	case bin.Model:
		m.binModel = v
	}
	return m, cmd
}

// handleSelectState manages the item type selection screen
func handleSelectState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, func() tea.Msg { return CancelMsg{} }
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case tea.KeyEnter:
			switch m.choices[m.cursor] {
			case i18n.AddPassword:
				m.state = AddPassword
			case i18n.AddCard:
				m.state = AddCard
			case i18n.AddText:
				m.state = AddText
			case i18n.AddBinary:
				m.state = AddBinary
			}
		}
	}
	return m, nil
}

// handleProcessingState manages the item addition operation
func handleProcessingState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case ErrorMsg:
		m.errMsg = msg.Err.Error()
		m.state = ErrorState
	case AddSuccessMsg:
		m.state = SelectState
	}
	return m, nil
}

// handleErrorState manages error display and recovery
func handleErrorState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.state = SelectState // Return to selection on confirmation
		}
	}
	return m, nil
}

// add performs the actual item storage operation
func (m Model) add(info *models.ItemInfo, content []byte) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
		defer cancel()

		err := m.service.Add(ctx, info, content)
		if err != nil {
			return ErrorMsg{err}
		}
		return AddSuccessMsg{}
	}
}
