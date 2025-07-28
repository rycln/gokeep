package add

import (
	"context"
	"time"

	"github.com/rycln/gokeep/internal/client/tui/items/card"
	"github.com/rycln/gokeep/internal/client/tui/items/logpass"
	"github.com/rycln/gokeep/internal/shared/models"
)

type state int

const (
	SelectState state = iota
	AddPassword
	AddCard
	ProcessingState
	ErrorState
)

type choice string

const (
	password = "Логин/Пароль"
	bankcard = "Банковская карта"
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
	cardModel    card.Model
	user         *models.User
	service      itemAdder
	timeout      time.Duration
}

func InitialModel(service itemAdder, timeout time.Duration) Model {
	return Model{
		state:        SelectState,
		choices:      []choice{password, bankcard, text, binary},
		logpassModel: logpass.InitialModel(),
		cardModel:    card.InitialModel(),
		service:      service,
		timeout:      timeout,
	}
}

func (m *Model) SetUser(user *models.User) {
	m.user = user
}
