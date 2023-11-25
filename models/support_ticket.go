package models

import (
	"time"
)

// SupportTicket represents a support ticket
type SupportTicket struct {
	ID        uint            `gorm:"primarykey" json:"-"`
	UserID    uint            `gorm:"not null" json:"-"`
	Username  string          `gorm:"not null" json:"username"`
	Subject   string          `gorm:"type:varchar(100);not null" json:"subject"`
	Messages  []TicketMessage `json:"messages"`
	TradeId   *uint           `json:"trade_id,omitempty"`
	Status    int             `gorm:"default:0;not null" json:"status"`
	CreatedAt time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time       `gorm:"not null" json:"updated_at"`
}

// TicketMessage represents a message within a support ticket
type TicketMessage struct {
	ID              uint   `gorm:"primarykey" json:"-"`
	SupportTicketID uint   `json:"-"`
	Msg             string `gorm:"type:text;not null" json:"message"`
	SenderUsername  string `gorm:"not null" json:"sender_username"`
	CreatedAt       time.Time
}

const (
	OpenTicket = iota
	PendingTicket
	ClosedTicket
)

func (SupportTicket) TableName() string {
	return "support_tickets"
}

func (TicketMessage) TableName() string {
	return "ticket_messages"
}
