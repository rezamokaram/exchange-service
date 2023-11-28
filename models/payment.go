package models

import (
	"time"

	"gorm.io/gorm"
)

type PaymentInfo struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID    uint   `gorm:"not null"`
	Amount    int64  `gorm:"type:bigint"`
	Status    string `gorm:"type:varchar(255)"`
	Authority string `gorm:"type:varchar(255)"`
}

func (PaymentInfo) TableName() string {
	return "payment_info"
}
