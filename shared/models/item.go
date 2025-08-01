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
	ID        ItemID    // Unique item identifier
	UserID    UserID    // ID of the user who owns the item
	ItemType  ItemType  // Type of the item (password, card, etc.)
	Name      string    // Human-readable name of the item
	Metadata  string    // Additional metadata in string format
	UpdatedAt time.Time // Last modification timestamp
}

// Item represents a complete stored item with all data fields.
// Used for full item operations including synchronization.
type Item struct {
	ID        ItemID    // Unique item identifier
	UserID    UserID    // ID of the user who owns the item
	ItemType  ItemType  // Type of the item (password, card, etc.)
	Name      string    // Human-readable name of the item
	Metadata  string    // Additional metadata in string format
	Data      []byte    // Encrypted item data
	UpdatedAt time.Time // Last modification timestamp
	IsDeleted bool      // Soft delete flag
}
