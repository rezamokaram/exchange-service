package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username      string `gorm:"type:varchar(255);unique;not null"`
	Email         string `gorm:"type:varchar(255);unique;not null"`
	Password      string `gorm:"type:varchar(255);not null"`
	IsAdmin       bool   `gorm:"default:false"`
	Profile       Profile
	Transaction   []Transaction
	BankingInfo   []BankingInfo
	SupportTicket []SupportTicket
}

func (User) TableName() string {
	return "users"
}

// UserInfo represents a user's detailed information in the system
type UserInfo struct {
	Username            string   `json:"username" example:"john_doe"`                // User's username
	Email               string   `json:"email" example:"john@example.com"`           // User's email address
	IsAdmin             bool     `json:"is_admin" example:"true"`                    // Indicates if user has admin privileges
	PhoneNumber         string   `json:"phone_number" example:"123-456-7890"`        // User's phone number
	AuthenticationLevel int      `json:"authentication_level" example:"2"`           // User's authentication level
	BlockedLevel        int      `json:"blocked_level" example:"0"`                  // User's blocked level
	Balance             int      `json:"balance" example:"1000"`                     // User's account balance
	IsPremium           bool     `json:"is_premium" example:"false"`                 // Indicates if user has a premium account
	BanksNames          []string `json:"banks_names" example:"['Bank A', 'Bank B']"` // List of user's bank names
}

func NewUserInfo(user User) UserInfo {
	newUserInfo := UserInfo{}
	newUserInfo.Username = user.Username
	newUserInfo.Email = user.Email
	newUserInfo.PhoneNumber = user.Profile.PhoneNumber
	newUserInfo.AuthenticationLevel = user.Profile.AuthenticationLevel
	newUserInfo.BlockedLevel = user.Profile.BlockedLevel
	newUserInfo.Balance = user.Profile.Balance
	newUserInfo.IsPremium = user.Profile.IsPremium

	for _, bi := range user.BankingInfo {
		newUserInfo.BanksNames = append(newUserInfo.BanksNames, bi.BankName)
	}

	return newUserInfo
}
