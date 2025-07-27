package models

type Card struct {
	CardNumber string `json:"card_number"`
	CardOwner  string `json:"card_owner"`
	ExpiryDate string `json:"expiry_date"`
	CVV        string `json:"cvv"`
}
