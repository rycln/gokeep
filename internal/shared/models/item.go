package models

import (
	"fmt"
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

func (i ItemInfo) FilterValue() string { return i.Name }
func (i ItemInfo) Title() string       { return i.Name }
func (i ItemInfo) Description() string { return fmt.Sprintf("Type: %s", i.ItemType) }
