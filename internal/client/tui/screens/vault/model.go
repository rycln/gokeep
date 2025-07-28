package vault

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/rycln/gokeep/internal/shared/models"
)

type state int

const (
	ListState state = iota
	DetailState
	ProcessingState
	ErrorState
)

type (
	AddItemReqMsg    struct{ User *models.User }
	GetContentReqMsg struct{}
	ItemsMsg         struct{ Items []itemRender }
	ErrorMsg         struct{ Err error }
)

type itemGetter interface {
	List(context.Context, models.UserID) ([]models.ItemInfo, error)
	GetContent(context.Context, string) ([]byte, error)
}

type itemRender struct {
	ItemType  models.ItemType
	Name      string
	Metadata  string
	UpdatedAt time.Time
}

func (i itemRender) FilterValue() string { return i.Name }
func (i itemRender) Title() string       { return i.Name }
func (i itemRender) Description() string {
	return fmt.Sprintf("Тип: %s\n Описание: %s", i.ItemType, i.Metadata)
}

type Model struct {
	state    state
	selected *itemRender
	items    []itemRender
	list     list.Model
	errMsg   string
	service  itemGetter
	user     *models.User
	timeout  time.Duration
}

func InitialModel(service itemGetter, timeout time.Duration) Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 20, 20)
	l.Title = "Goph Keeper"
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	return Model{
		state:   ListState,
		list:    l,
		service: service,
		timeout: timeout,
	}
}

func (m *Model) SetUser(user *models.User) {
	m.user = user
}
