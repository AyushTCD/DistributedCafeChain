package model

import (
	"time"
)

// User represents a user in the loyalty service system.
type User struct {
	ID           string  `gorm:"column:user_uuid"`
	AccountID    *string `gorm:"column:account_uuid"`
	Account      *Account
	Name         string
	Email        string `gorm:"unique;column:email_address"`
	Password     string
	Phone        string    `gorm:"unique;column:phone_number"`
	CreationDate time.Time `gorm:"autoCreateTime"`
	InviteCode   *string   `gorm:"uniqueIndex;"`
}
