package model

import (
	"time"
)

// Transaction represents a transaction in the loyalty service system.
type Transaction struct {
	ID           string    `gorm:"column:transaction_uuid"` // unique ID
	AccountID    string    `gorm:"column:account_uuid"`
	UserID       string    `gorm:"column:user_uuid"` // ID of the user who made the transaction
	Amount       float64   // Transaction amount
	Date         time.Time `gorm:"autoCreateTime"`
	PointsEarned int
}
