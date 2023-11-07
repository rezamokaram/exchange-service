package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	UserID              uint   `gorm:"not null"`
	Email               string `gorm:"type:varchar(100);not null"`
	PhoneNumber         string `gorm:"type:varchar(100)"`
	AuthenticationLevel int    `gorm:"default:0;not null"`
	BlockedLevel        int    `gorm:"default:0;not null"`
	Balance             int    `gorm:"default:0;not null"`
	IsPremium           bool   `gorm:"default:false"`
}

func (Profile) TableName() string {
	return "profiles"
}
