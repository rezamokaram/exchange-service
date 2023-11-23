package models

import "gorm.io/gorm"

type SupportTicket struct {
	gorm.Model
	UserID      uint   `gorm:"not null"`
	Subject     string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:text"`
	TradeId     *uint
	Status      int `gorm:"default:0;not null"`
}

func (SupportTicket) TableName() string {
	return "support_tickets"
}
