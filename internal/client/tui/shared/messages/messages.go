package messages

import "github.com/rycln/gokeep/internal/shared/models"

type (
	ItemMsg struct {
		Info    *models.ItemInfo
		Content []byte
	}
	ErrMsg    struct{ Err error }
	CancelMsg struct{}
)
