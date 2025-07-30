// Card package manages credit card data storage and display.
package card

// Card represents payment card information for secure storage.
// Contains sensitive payment details that should be encrypted.
type Card struct {
	CardNumber string `json:"card_number"` // Full card number (16-19 digits)
	CardOwner  string `json:"card_owner"`  // Name of cardholder
	ExpiryDate string `json:"expiry_date"` // MM/YY expiration format
	CVV        string `json:"cvv"`         // 3-4 digit security code
}
