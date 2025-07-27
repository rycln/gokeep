package vault

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
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
	var cmd tea.Cmd

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
			m.state = ListState
			m.selected = nil
		case tea.KeyEnter:
			return m, func() tea.Msg { return GetContentReqMsg{} }
		}
	}

	return m, nil
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
	}

	return m, nil
}
