package models

import "time"

// ItemID represents a unique identifier for stored items.
type ItemID string

// ItemType categorizes stored items into specific types.
type ItemType string

// Supported item type constants
const (
	TypePassword ItemType = "Пароль"
	TypeCard     ItemType = "Карта"
	TypeText     ItemType = "Текст"
	TypeBinary   ItemType = "Бинарный файл"
)

// ItemInfo contains metadata about a stored item.
// Represents common fields across all item types.
type ItemInfo struct {
	ID        ItemID
	UserID    UserID
	ItemType  ItemType
	Name      string
	Metadata  string
	UpdatedAt time.Time
}
