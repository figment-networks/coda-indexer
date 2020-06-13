package model

import (
	"errors"
	"fmt"
	"time"
)

const (
	TransactionTypePayment    = "payment"
	TransactionTypeDelegation = "delegation"
)

// Transaction contains the blockchain transaction details
type Transaction struct {
	Model

	Type      string    `json:"type"`
	Hash      string    `json:"hash"`
	BlockHash string    `json:"block_hash"`
	Height    uint64    `json:"height"`
	Time      time.Time `json:"time"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Amount    uint64    `json:"amount"`
	Fee       uint64    `json:"fee"`
	Nonce     int       `json:"nonce"`
	Memo      string    `json:"memo"`
}

type TransactionCount struct {
	Time  time.Time
	Count int
}

// TableName returns the model table name
func (Transaction) TableName() string {
	return "transactions"
}

// String returns transaction text representation
func (t Transaction) String() string {
	return fmt.Sprintf("type=%v hash=%v height=%v", t.Type, t.Hash, t.Height)
}

// Validate returns an error if transaction is invalid
func (t Transaction) Validate() error {
	if t.Type == "" {
		return errors.New("type is required")
	}
	if t.BlockHash == "" {
		return errors.New("block hash is required")
	}
	if t.Hash == "" {
		return errors.New("hash is required")
	}
	if t.Height <= 0 {
		return errors.New("height is invalid")
	}
	if t.Time.IsZero() {
		return errors.New("time is invalid")
	}
	if t.Sender == "" {
		return errors.New("sender is required")
	}
	if t.Receiver == "" {
		return errors.New("receiver is required")
	}
	if t.Amount < 0 {
		return errors.New("amount is invalid")
	}
	return nil
}
