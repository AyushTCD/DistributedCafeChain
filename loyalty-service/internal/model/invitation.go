package model

import "time"

type Invitation struct {
	InvitationUUID string    `gorm:"primaryKey;column:invitation_uuid"`
	Email          string    `gorm:"not null;column:email"`
	AccountUUID    string    `gorm:"column:account_uuid"`
	InviterUUID    string    `gorm:"column:inviter_uuid"`
	Token          string    `gorm:"unique;not null;column:token"`
	CreationDate   time.Time `gorm:"not null;column:creation_date"`
	ExpirationDate time.Time `gorm:"not null;column:expiration_date"`
	Status         string    `gorm:"not null;column:status"`
}
