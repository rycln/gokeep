package update

import (
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/internal/client/tui/items/bin"
	"github.com/rycln/gokeep/internal/client/tui/items/card"
	"github.com/rycln/gokeep/internal/client/tui/items/logpass"
	"github.com/rycln/gokeep/internal/client/tui/items/text"
	"github.com/rycln/gokeep/internal/client/tui/shared/messages"
	"github.com/rycln/gokeep/internal/shared/models"
)

var errEmptyItem = errors.New("no info or content")

func (m Model) Init() tea.Cmd {
	return nil
}

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
			updated, cmd := m.logpassModel.Update(msg)
			if logpassModel, ok := updated.(logpass.Model); ok {
				m.logpassModel = logpassModel
			}
			return m, cmd
		case UpdateCard:
			updated, cmd := m.cardModel.Update(msg)
			if cardModel, ok := updated.(card.Model); ok {
				m.cardModel = cardModel
			}
			return m, cmd
		case UpdateText:
			updated, cmd := m.textModel.Update(msg)
			if textModel, ok := updated.(text.Model); ok {
				m.textModel = textModel
			}
			return m, cmd
		case UpdateBinary:
			updated, cmd := m.binModel.Update(msg)
			if binModel, ok := updated.(bin.Model); ok {
				m.binModel = binModel
			}
			return m, cmd
		case ErrorState:
			return handleErrorState(m, msg)
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

func (m Model) initUpdate() tea.Cmd {
	return func() tea.Msg {
		if m.info == nil || m.content == nil {
			return ErrorMsg{errEmptyItem}
		}

		var s state
		switch m.info.ItemType {
		case models.TypePassword:
			err := m.logpassModel.SetStartData(m.info, m.content)
			if err != nil {
				return ErrorMsg{err}
			}
			s = UpdatePassword
		case models.TypeCard:

		case models.TypeText:

		case models.TypeBinary:

		}

		return InitSuccessMsg{state: s}
	}
}

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
		return m, nil
	case InitSuccessMsg:
		m.state = msg.state
		return m, nil
	case UpdateSuccessMsg:
		return m, func() tea.Msg { return CancelMsg{} }
	}

	return m, nil
}

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
