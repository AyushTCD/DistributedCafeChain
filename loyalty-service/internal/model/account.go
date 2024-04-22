package model

import (
	"time"
)

// Account represents a loyalty group account
type Account struct {
	ID           string `gorm:"column:account_uuid"`
	Users        []User
	Points       int       `gorm:"column:points_balance"`
	CreationDate time.Time `gorm:"autoCreateTime"`
}
