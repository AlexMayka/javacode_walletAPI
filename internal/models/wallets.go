// Package models contains data structures used across the wallet service.
package models

import "time"

// Wallet represents a user's wallet with balance and timestamps.
type Wallet struct {
	Id          string
	Balance     uint64
	CreatedTime time.Time
	UpdatedTime time.Time
}

// WalletOperationRequest represents the request body for a wallet operation
// such as deposit or withdrawal.
type WalletOperationRequest struct {
	// WalletID is the unique identifier of the wallet.
	// required: true
	// example: abc123
	WalletID string `json:"walletId" example:"c3a8cb84-03f2-4fb9-982a-9ee2cfb50b9f"`

	// OperationType defines the type of operation: "deposit" or "withdrawal".
	// required: true
	// example: deposit
	OperationType string `json:"operationType" example:"DEPOSIT"`

	// Amount is the amount of money to deposit or withdraw.
	// Must be a positive integer.
	// required: true
	// example: 500
	Amount int `json:"amount" example:"1000"`
}

// BalanceResponse represents the response containing the wallet balance.
type BalanceResponse struct {
	Uuid    string `json:"uuid" example:"c3a8cb84-03f2-4fb9-982a-9ee2cfb50b9f"`
	Balance uint64 `json:"balance" example:"1000"`
}
