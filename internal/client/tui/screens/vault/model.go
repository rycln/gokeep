package vault

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/rycln/gokeep/internal/client/tui/shared/i18n"
	"github.com/rycln/gokeep/internal/client/tui/shared/styles"
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
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("62")).
		Foreground(lipgloss.Color("229"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("245"))
	delegate.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("n"),
				key.WithHelp("n", "добавить"),
			),
			key.NewBinding(
				key.WithKeys("u"),
				key.WithHelp("u", "обновить"),
			),
		}
	}

	l := list.New([]list.Item{}, delegate, 20, 20)
	l.Title = i18n.VaultTitle
	l.Styles.Title = styles.TitleStyle
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	l.StatusMessageLifetime = timeout

	l.SetStatusBarItemName(i18n.VaultListTitleNameSingular, i18n.VaultListTitleNamePlural)

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
