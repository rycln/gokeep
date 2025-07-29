package vault

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/internal/client/tui/items/card"
	"github.com/rycln/gokeep/internal/client/tui/items/logpass"
	"github.com/rycln/gokeep/internal/client/tui/items/text"
	"github.com/rycln/gokeep/internal/shared/models"
)

func (m Model) Init() tea.Cmd {
	if m.user == nil {
		return func() tea.Msg {
			return ErrorMsg{Err: fmt.Errorf("user not set")}
		}
	}

	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case UpdateState:
		m.state = ProcessingState
		return m, m.loadItems()
	case ListState:
		return handleListState(m, msg)
	case DetailState:
		return handleDetailState(m, msg)
	case ProcessingState:
		return handleProcessingState(m, msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

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
			case "u":
				m.state = ProcessingState
				return m, m.loadItems()
			case "n":
				return m, func() tea.Msg { return AddItemReqMsg{User: m.user} }
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

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

func handleDetailState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc, tea.KeyBackspace:
			m.state = UpdateState
			m.selected = nil
		case tea.KeyEnter:
			m.state = ProcessingState
			return m, m.getContent()
		}
	}

	return m, nil
}

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
			if err != nil {
				return ErrorMsg{Err: err}
			}
		case models.TypeCard:
			content, err = card.GetContentRender(contentBytes)
			if err != nil {
				return ErrorMsg{Err: err}
			}
		case models.TypeText:
			content, err = text.GetContentRender(contentBytes)
			if err != nil {
				return ErrorMsg{Err: err}
			}
		}

		return ContentMsg{Content: content}
	}
}

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
	}

	return m, nil
}
