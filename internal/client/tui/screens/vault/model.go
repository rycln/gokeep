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
	UpdateState
	BinaryInputState
	ProcessingState
	ErrorState
)

type (
	AddItemReqMsg struct{ User *models.User }
	UpdateReqMsg  struct {
		Info    *models.ItemInfo
		Content []byte
	}
	DeleteSuccessMsg struct{}
	ItemsMsg         struct{ Items []itemRender }
	ContentMsg       struct{ Content string }
	ErrorMsg         struct{ Err error }
)

type itemGetter interface {
	List(context.Context, models.UserID) ([]models.ItemInfo, error)
	GetContent(context.Context, models.ItemID) ([]byte, error)
}

type itemDeleter interface {
	Delete(context.Context, models.ItemID) error
}

type itemService interface {
	itemGetter
	itemDeleter
}

type itemRender struct {
	ID        models.ItemID
	ItemType  models.ItemType
	Name      string
	Metadata  string
	UpdatedAt time.Time
	Content   string
}

func (i itemRender) FilterValue() string { return i.Name }
func (i itemRender) Title() string       { return i.Name }
func (i itemRender) Description() string {
	return fmt.Sprintf("Тип: %s\n Описание: %s", i.ItemType, i.Metadata)
}

type Model struct {
	state    state
	input    string
	selected *itemRender
	items    []itemRender
	list     list.Model
	errMsg   string
	service  itemService
	user     *models.User
	timeout  time.Duration
}

func InitialModel(service itemService, timeout time.Duration) Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 20, 20)
	l.Title = "Goph Keeper"
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	return Model{
		state:   UpdateState,
		list:    l,
		service: service,
		timeout: timeout,
	}
}

func (m *Model) SetUser(user *models.User) {
	m.user = user
}

func (m *Model) SetUpdateState() {
	m.state = UpdateState
}
