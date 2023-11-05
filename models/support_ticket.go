package models

import "gorm.io/gorm"

type SupportTicket struct {
	gorm.Model
	UserID        uint   `gorm:"not null"`
	Subject       string `gorm:"type:varchar(100);not null"`
	Description   string `gorm:"type:text"`
	TransactionId uint   `gorm:""`
	Status        int    `gorm:"default:0;not null"` /* 0: not read, 1: reading, 2: read */
}

func (SupportTicket) TableName() string {
	return "support_tickets"
}
