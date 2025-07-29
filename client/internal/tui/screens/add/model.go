package add

import (
	"context"
	"time"

	"github.com/rycln/gokeep/client/internal/tui/items/bin"
	"github.com/rycln/gokeep/client/internal/tui/items/card"
	"github.com/rycln/gokeep/client/internal/tui/items/logpass"
	"github.com/rycln/gokeep/client/internal/tui/items/text"
	"github.com/rycln/gokeep/shared/models"
)

type state int

const (
	SelectState state = iota
	AddPassword
	AddCard
	AddText
	AddBinary
	ProcessingState
	ErrorState
)

type choice string

const (
	password    = "Логин/Пароль"
	bankcard    = "Банковская карта"
	textcontent = "Текст"
	binary      = "Бинарный файл"
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
	textModel    text.Model
	binModel     bin.Model
	user         *models.User
	service      itemAdder
	timeout      time.Duration
}

func InitialModel(service itemAdder, timeout time.Duration) Model {
	return Model{
		state:        SelectState,
		choices:      []choice{password, bankcard, textcontent, binary},
		logpassModel: logpass.InitialModel(),
		cardModel:    card.InitialModel(),
		textModel:    text.InitialModel(),
		binModel:     bin.InitialModel(),
		service:      service,
		timeout:      timeout,
	}
}

func (m *Model) SetUser(user *models.User) {
	m.user = user
}
