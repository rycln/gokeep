package models

import (
	"time"
)

type ItemID string

type ItemType string

const (
	TypePassword ItemType = "password"
	TypeCard     ItemType = "card"
	TypeText     ItemType = "text"
	TypeBinary   ItemType = "binary"
)

type ItemInfo struct {
	ID        ItemID
	UserID    UserID
	ItemType  ItemType
	Name      string
	Metadata  string
	UpdatedAt time.Time
}
