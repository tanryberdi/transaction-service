package entities

import "github.com/shopspring/decimal"

type SourceType string

// User represents a user in the system
type User struct {
	ID      uint64          `json:"id" db:"id"`
	Balance decimal.Decimal `json:"balance" db:"balance"`
}

// Transaction represents a transaction in the system
type Transaction struct {
	ID            uint64          `json:"id" db:"id"`
	UserID        uint64          `json:"userId" db:"user_id"`
	TransactionID string          `json:"transactionId" db:"transaction_id"`
	State         string          `json:"state" db:"state"`
	Amount        decimal.Decimal `json:"amount" db:"amount"`
	SourceType    SourceType      `json:"sourceType" db:"source_type"`
	CreatedAt     string          `json:"createdAt" db:"created_at"`
}

// TransactionState represents the state of a transaction
type TransactionState string

const (
	StateWin  TransactionState = "win"
	StateLose TransactionState = "lose"
)

// IsValid checks if the transaction state is valid
func (ts TransactionState) IsValid() bool {
	return ts == StateWin || ts == StateLose
}

const (
	SourceTypeGame    SourceType = "game"
	SourceTypeServer  SourceType = "server"
	SourceTypePayment SourceType = "payment"
)

// IsValid checks if the source type is valid
func (st SourceType) IsValid() bool {
	return st == SourceTypeGame || st == SourceTypeServer || st == SourceTypePayment
}

// TransactionRequest represents the incoming transaction request
type TransactionRequest struct {
	State         string `json:"state" binding:"required"`
	Amount        string `json:"amount" binding:"required"`
	TransactionID string `json:"transactionId" binding:"required"`
}

type BalanceResponse struct {
	UserID  uint64 `json:"userId"`
	Balance string `json:"balance"`
}
