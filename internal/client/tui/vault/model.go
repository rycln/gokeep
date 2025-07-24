package vault

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/internal/shared/models"
)

type itemAdder interface {
	Add(context.Context, *models.ItemInfo, []byte) error
}

type itemGetter interface {
	List(context.Context, models.UserID) ([]models.ItemInfo, error)
	GetContent(context.Context, string) ([]byte, error)
}

type itemService interface {
	itemAdder
	itemGetter
}

type Model struct {
	list     list.Model
	items    []models.ItemInfo
	selected *models.ItemInfo
	content  string
	service  itemService
	ctx      context.Context
	mode     string // "list" or "detail"
}

type Msg struct {
	Items   []models.ItemInfo
	Content []byte
	Err     error
}

func NewModel(service itemService, ctx context.Context) Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Your Vault"
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return Model{
		list:    l,
		service: service,
		ctx:     ctx,
		mode:    "list",
	}
}

func (m Model) LoadItems(userID models.UserID) tea.Cmd {
	return func() tea.Msg {
		items, err := m.service.List(m.ctx, userID)
		if err != nil {
			return Msg{Err: err}
		}
		return Msg{Items: items}
	}
}

func (m Model) LoadItemContent(name string) tea.Cmd {
	return func() tea.Msg {
		content, err := m.service.GetContent(m.ctx, name)
		if err != nil {
			return Msg{Err: err}
		}
		return Msg{Content: content}
	}
}
