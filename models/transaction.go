package models

import "gorm.io/gorm"

// Transaction represents a cryptocurrency transaction
type Transaction struct {
	gorm.Model
	UserID      uint   `gorm:"not null" json:"user_id,omitempty" example:"2"`
	Amount      int    `gorm:"not null" json:"amount,omitempty" example:"515"`
	Method      bool   `gorm:"not null" json:"method,omitempty" example:"0"`  // method: false -> withdraw, true -> deposit
	Service     int    `gorm:"not null" json:"service,omitempty" example:"1"` // service: 0 -> payment, 1 -> trade
	Description string `gorm:"" json:"description,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}

func NewTransaction(userID uint, amount, service int, method bool, description string) Transaction {
	return Transaction{
		UserID: userID,
		Amount: amount,
		Service: service,
		Method: method,
		Description: description,
	}
}
