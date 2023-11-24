package models

import "gorm.io/gorm"

type SupportTicket struct {
	gorm.Model  `json:"-"`
	UserID      uint   `gorm:"not null" json:"-"`
	Username    string `gorm:"not null" json:"username"`
	Subject     string `gorm:"type:varchar(100);not null" json:"subject"`
	Description string `gorm:"type:text" json:"description"`
	TradeId     *uint  `json:"trade_id,omitempty"`
	Status      int    `gorm:"default:0;not null" json:"status,omitempty"`
}

func (SupportTicket) TableName() string {
	return "support_tickets"
}
