package models

import (
	"gorm.io/gorm"
)

type PaymentInfo struct {
	gorm.Model
	UserID    uint   `gorm:"not null"`
	Amount    int64  `gorm:"type:bigint"`
	Status    string `gorm:"type:varchar(255)"`
	Authority string `gorm:"type:varchar(255)"`
}

func (PaymentInfo) TableName() string {
	return "payment_info"
}
