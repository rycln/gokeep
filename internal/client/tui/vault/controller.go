package vault

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/internal/shared/models"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case Msg:
		return m.handleMsg(msg)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch m.mode {
	case "list":
		switch msg.String() {
		case "enter":
			if len(m.list.Items()) > 0 {
				selected := m.list.SelectedItem().(models.ItemInfo)
				m.selected = &selected
				m.mode = "detail"
				return m, m.LoadItemContent(selected.Name)
			}
		case "n":
			// TODO: Add new item
		}
	case "detail":
		switch msg.String() {
		case "esc", "backspace":
			m.mode = "list"
			m.selected = nil
			m.content = ""
		}
	}
	return m, nil
}

func (m Model) handleMsg(msg Msg) (Model, tea.Cmd) {
	if msg.Err != nil {
		// Handle error
		return m, nil
	}
	if msg.Items != nil {
		m.items = msg.Items
		items := make([]list.Item, len(msg.Items))
		for i, item := range msg.Items {
			items[i] = item
		}
		m.list.SetItems(items)
	}
	if msg.Content != nil {
		m.content = string(msg.Content)
	}
	return m, nil
}
