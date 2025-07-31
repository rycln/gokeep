package vault

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/client/internal/tui/items/bin"
	"github.com/rycln/gokeep/client/internal/tui/items/card"
	"github.com/rycln/gokeep/client/internal/tui/items/logpass"
	"github.com/rycln/gokeep/client/internal/tui/items/text"
	"github.com/rycln/gokeep/shared/models"
)

// Init initializes the vault model and checks for required user
func (m Model) Init() tea.Cmd {
	if m.user == nil {
		return func() tea.Msg {
			return ErrorMsg{Err: fmt.Errorf("user not set")}
		}
	}
	return nil
}

// Update handles all messages and state transitions
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case UpdateState:
		m.state = ProcessingState
		return m, m.loadItems()
	case StartState:
		m.state = ProcessingState
		return m, m.setKey()
	case ListState:
		return handleListState(m, msg)
	case DetailState:
		return handleDetailState(m, msg)
	case ProcessingState:
		return handleProcessingState(m, msg)
	case ErrorState:
		return handleErrorState(m, msg)
	case BinaryInputState:
		return handleBinaryInputState(m, msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// handleListState manages the item list view interactions
func handleListState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			if len(m.list.Items()) > 0 {
				selected := m.list.SelectedItem().(itemRender)
				m.selected = &selected
				m.state = DetailState
				return m, nil
			}
		case tea.KeyRunes:
			switch msg.String() {
			case "q", "й":
				return m, tea.Quit
			case "u", "г":
				m.state = ProcessingState
				return m, m.loadItems()
			case "n", "т":
				return m, func() tea.Msg { return AddItemReqMsg{User: m.user} }
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// loadItems fetches items from service for current user
func (m Model) loadItems() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
		defer cancel()

		items, err := m.service.List(ctx, m.user.ID)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		ritems := make([]itemRender, len(items))
		for i, item := range items {
			ritems[i] = itemRender{
				ID:        item.ID,
				ItemType:  item.ItemType,
				Name:      item.Name,
				Metadata:  item.Metadata,
				UpdatedAt: item.UpdatedAt,
			}
		}

		return ItemsMsg{Items: ritems}
	}
}

// handleDetailState manages item detail view interactions
func handleDetailState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc, tea.KeyBackspace:
			m.state = ListState
			m.selected = nil
		case tea.KeyEnter:
			if m.selected.ItemType == models.TypeBinary {
				m.state = BinaryInputState
				m.input = ""
				return m, nil
			}
			m.state = ProcessingState
			return m, m.getContent()
		case tea.KeyDelete:
			m.state = ProcessingState
			return m, m.deleteItem()
		case tea.KeyInsert:
			return m, m.updateItem()
		}
	}

	return m, nil
}

// getContent retrieves and formats item content based on type
func (m Model) getContent() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
		defer cancel()

		contentBytes, err := m.service.GetContent(ctx, m.selected.ID)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		var content string
		switch m.selected.ItemType {
		case models.TypePassword:
			content, err = logpass.GetContentRender(contentBytes)
		case models.TypeCard:
			content, err = card.GetContentRender(contentBytes)
		case models.TypeText:
			content, err = text.GetContentRender(contentBytes)
		case models.TypeBinary:
			content, err = bin.UploadFile(m.input, contentBytes)
		}
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return ContentMsg{Content: content}
	}
}

// deleteItem removes the currently selected item
func (m Model) deleteItem() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
		defer cancel()

		err := m.service.Delete(ctx, m.selected.ID)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return DeleteSuccessMsg{}
	}
}

// updateItem prepares item data for update screen
func (m Model) updateItem() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
		defer cancel()

		contentBytes, err := m.service.GetContent(ctx, m.selected.ID)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return UpdateReqMsg{
			Info: &models.ItemInfo{
				ID:        m.selected.ID,
				UserID:    m.user.ID,
				ItemType:  m.selected.ItemType,
				Name:      m.selected.Name,
				Metadata:  m.selected.Metadata,
				UpdatedAt: m.selected.UpdatedAt,
			},
			Content: contentBytes,
		}
	}
}

// handleProcessingState processes background operation results
func handleProcessingState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ErrorMsg:
		m.errMsg = msg.Err.Error()
		m.state = ErrorState
	case ItemsMsg:
		m.items = msg.Items
		m.state = ListState

		items := make([]list.Item, len(msg.Items))
		for i, item := range msg.Items {
			items[i] = item
		}
		return m, m.list.SetItems(items)
	case ContentMsg:
		m.selected.Content = msg.Content
		m.state = DetailState
	case DeleteSuccessMsg:
		m.state = UpdateState
	}

	return m, nil
}

// handleBinaryInputState manages binary file path input
func handleBinaryInputState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.state = ProcessingState
			return m, m.getContent()
		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case tea.KeyRunes:
			m.input += msg.String()
		}
	}

	return m, nil
}

// handleErrorState manages error display and recovery
func handleErrorState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.state = ListState
		}
	}

	return m, nil
}

// setKey set key into crypt
func (m Model) setKey() tea.Cmd {
	return func() tea.Msg {

		err := m.crypt.SetKey(m.user)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return DeleteSuccessMsg{}
	}
}
