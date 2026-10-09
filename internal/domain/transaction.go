package domain

import "time"

type Direction int

const (
	DirectionUnknown Direction = iota
	DirectionIn
	DirectionOut
)

type Transaction struct {
	ExternalID string
	AccountID  int64
	Date       time.Time
	Amount     float64
	Currency   string
	Direction  Direction

	Description string
	Merchant    string

	// shared by the two sides of a move between the user's own accounts
	TransferRef string

	// non-empty triggers a linked receipt after the tx is posted
	ReceiptItems []ReceiptItem
}

type ReceiptItem struct {
	Name        string
	AmountCents int64
}
