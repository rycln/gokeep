package models

import (
	"time"
)

type ItemID string

type ItemType string

const (
	TypePassword ItemType = "Пароль"
	TypeCard     ItemType = "Карта"
	TypeText     ItemType = "Текст"
	TypeBinary   ItemType = "Бинарный файл"
)

type ItemInfo struct {
	ID        ItemID
	UserID    UserID
	ItemType  ItemType
	Name      string
	Metadata  string
	UpdatedAt time.Time
}
