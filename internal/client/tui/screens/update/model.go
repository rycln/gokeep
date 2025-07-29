package update

import (
	"context"
	"time"

	"github.com/rycln/gokeep/internal/client/tui/items/bin"
	"github.com/rycln/gokeep/internal/client/tui/items/card"
	"github.com/rycln/gokeep/internal/client/tui/items/logpass"
	"github.com/rycln/gokeep/internal/client/tui/items/text"
	"github.com/rycln/gokeep/internal/shared/models"
)

type state int

const (
	LoadState state = iota
	UpdatePassword
	UpdateCard
	UpdateText
	UpdateBinary
	ProcessingState
	ErrorState
)

type itemUpdater interface {
	Update(context.Context, *models.ItemInfo, []byte) error
}

type (
	InitSuccessMsg   struct{ state state }
	UpdateSuccessMsg struct{}
	ErrorMsg         struct{ Err error }
	CancelMsg        struct{}
)

type Model struct {
	state        state
	errMsg       string
	logpassModel logpass.Model
	cardModel    card.Model
	textModel    text.Model
	binModel     bin.Model
	info         *models.ItemInfo
	content      []byte
	service      itemUpdater
	timeout      time.Duration
}

func InitialModel(service itemUpdater, timeout time.Duration) Model {
	return Model{
		state:        LoadState,
		logpassModel: logpass.InitialModel(),
		cardModel:    card.InitialModel(),
		textModel:    text.InitialModel(),
		binModel:     bin.InitialModel(),
		service:      service,
		timeout:      timeout,
	}
}

func (m *Model) SetItem(info *models.ItemInfo, content []byte) {
	m.state = LoadState
	m.info = info
	m.content = content
}
