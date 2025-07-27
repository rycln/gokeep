package add

import (
	"context"
	"time"

	"github.com/rycln/gokeep/internal/client/tui/add/logpass"
	"github.com/rycln/gokeep/internal/shared/models"
)

type state int

const (
	SelectState state = iota
	AddPassword
	ProcessingState
	ErrorState
)

type choice string

const (
	password = "Логин/Пароль"
	card     = "Банковская карта"
	text     = "Текст"
	binary   = "Бинарный файл"
)

type itemAdder interface {
	Add(context.Context, *models.ItemInfo, []byte) error
}

type (
	AddSuccessMsg struct{}
	ErrorMsg      struct{ Err error }
	CancelMsg     struct{}
)

type Model struct {
	state        state
	choices      []choice
	cursor       int
	errMsg       string
	logpassModel logpass.Model
	user         *models.User
	service      itemAdder
	timeout      time.Duration
}

func InitialModel(service itemAdder, timeout time.Duration) Model {
	return Model{
		state:        SelectState,
		choices:      []choice{password, card, text, binary},
		logpassModel: logpass.InitialModel(),
		service:      service,
		timeout:      timeout,
	}
}

func (m *Model) SetUser(user *models.User) {
	m.user = user
}
