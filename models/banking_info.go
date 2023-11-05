package models

import (
	"time"

	"gorm.io/gorm"
)

type BankingInfo struct {
	gorm.Model
	UserID        string    `gorm:"not null"`
	BankName      string    `gorm:"type:varchar(100);not null"`
	AccountNumber string    `gorm:"type:varchar(100);not null"`
	CardNumber    string    `gorm:"type:varchar(100);not null"`
	ExpireDate    time.Time `gorm:"not null"` /* => needs testing */
	Cvv2          string    `gorm:"type:varchar(5);not null"`
}

func (BankingInfo) TableName() string {
	return "banking_info"
}
