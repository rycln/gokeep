package models

import "time"

type ItemID string

type ItemType string

const (
	TypePassword ItemType = "password"
	TypeCard     ItemType = "card"
	TypeText     ItemType = "text"
	TypeBinary   ItemType = "binary"
)

type Item struct {
	ID        ItemID
	UserID    UserID
	ItemType  ItemType
	Name      string
	Metadata  string
	Content   []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}
