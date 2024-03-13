package user

import (
	"qexchange/models"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username      string `gorm:"type:varchar(255);unique;not null"`
	Email         string `gorm:"type:varchar(255);unique;not null"`
	Password      string `gorm:"type:varchar(255);not null"`
	IsAdmin       bool   `gorm:"default:false"`
	Profile       Profile
	Transaction   []models.Transaction
	BankingInfo   []models.BankingInfo
	SupportTicket []models.SupportTicket
}

func (User) TableName() string {
	return "users"
}
