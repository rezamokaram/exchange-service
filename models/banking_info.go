package models

import (
	"gorm.io/gorm"
)

type BankingInfo struct {
	gorm.Model
	UserID        string `gorm:"not null"`
	BankName      string `gorm:"type:varchar(100);not null" json:"bank_name"`
	AccountNumber string `gorm:"type:varchar(100);not null"`
	CardNumber    string `gorm:"type:varchar(100);not null"`
	ExpireDate    string `gorm:"not null"`
	Cvv2          string `gorm:"type:varchar(5);not null"`
}

func (BankingInfo) TableName() string {
	return "banking_info"
}
