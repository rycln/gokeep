package update

import (
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/client/internal/tui/items/bin"
	"github.com/rycln/gokeep/client/internal/tui/items/card"
	"github.com/rycln/gokeep/client/internal/tui/items/logpass"
	"github.com/rycln/gokeep/client/internal/tui/items/text"
	"github.com/rycln/gokeep/client/internal/tui/shared/messages"
	"github.com/rycln/gokeep/shared/models"
)

// errEmptyItem indicates missing item data
var errEmptyItem = errors.New("no info or content")

// Init initializes the update model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles all messages and state transitions for item updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.ItemMsg:
		m.state = ProcessingState
		return m, m.update(msg.Info, msg.Content)
	case messages.ErrMsg:
		m.errMsg = msg.Err.Error()
		m.state = ErrorState
		return m, nil
	case messages.CancelMsg:
		return m, func() tea.Msg { return CancelMsg{} }
	default:
		switch m.state {
		case UpdatePassword:
			return handleItemModelUpdate(m, msg, &m.logpassModel)
		case UpdateCard:
			return handleItemModelUpdate(m, msg, &m.cardModel)
		case UpdateText:
			return handleItemModelUpdate(m, msg, &m.textModel)
		case UpdateBinary:
			return handleItemModelUpdate(m, msg, &m.binModel)
		default:
			switch m.state {
			case LoadState:
				m.state = ProcessingState
				return m, m.initUpdate()
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

// update performs the actual item update operation
func (m Model) update(info *models.ItemInfo, content []byte) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
		defer cancel()

		info.ID = m.info.ID
		info.UserID = m.info.UserID

		err := m.service.Update(ctx, info, content)
		if err != nil {
			return ErrorMsg{err}
		}
		return UpdateSuccessMsg{}
	}
}

// initUpdate prepares the appropriate update form based on item type
func (m Model) initUpdate() tea.Cmd {
	return func() tea.Msg {
		if m.info == nil || m.content == nil {
			return ErrorMsg{errEmptyItem}
		}

		var s state
		var err error

		switch m.info.ItemType {
		case models.TypePassword:
			err = m.logpassModel.SetStartData(m.info, m.content)
			s = UpdatePassword
		case models.TypeCard:
			err = m.cardModel.SetStartData(m.info, m.content)
			s = UpdateCard
		case models.TypeText:
			err = m.textModel.SetStartData(m.info, m.content)
			s = UpdateText
		case models.TypeBinary:
			err = m.binModel.SetStartData(m.info, m.content)
			s = UpdateBinary
		}

		if err != nil {
			return ErrorMsg{err}
		}
		return InitSuccessMsg{state: s}
	}
}

// handleProcessingState manages the update operation status
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
	case InitSuccessMsg:
		m.state = msg.state
	case UpdateSuccessMsg:
		return m, func() tea.Msg { return CancelMsg{} }
	}
	return m, nil
}

// handleErrorState manages error display and recovery
func handleErrorState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, func() tea.Msg { return CancelMsg{} }
		}
	}
	return m, nil
}
