package models

import (
	"time"
)

// SupportTicket represents a support ticket
type SupportTicket struct {
	ID        uint            `gorm:"primarykey" json:"-"`
	UserID    uint            `gorm:"not null" json:"-"`
	Username  string          `gorm:"not null" json:"username" example:"user1"`
	Subject   string          `gorm:"type:varchar(100);not null" json:"subject" example:"Issue with Trade Execution"`
	Messages  []TicketMessage `json:"messages"`
	TradeId   *uint           `json:"trade_id,omitempty" example:"1"`
	Status    int             `gorm:"default:0;not null" json:"status" example:"0"`
	CreatedAt time.Time       `gorm:"not null" json:"created_at" example:"2023-01-01T15:04:05Z07:00"`
	UpdatedAt time.Time       `gorm:"not null" json:"updated_at" example:"2023-01-02T15:04:05Z07:00"`
}

// TicketMessage represents a message within a support ticket
type TicketMessage struct {
	ID              uint      `gorm:"primarykey" json:"-"`
	SupportTicketID uint      `json:"-"`
	Msg             string    `gorm:"type:text;not null" json:"message" example:"I encountered an error when trying to execute a trade."`
	SenderUsername  string    `gorm:"not null" json:"sender_username" example:"user1"`
	CreatedAt       time.Time `json:"created_at" example:"2023-01-01T15:04:05Z07:00"`
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
