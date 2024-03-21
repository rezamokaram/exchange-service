package user

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	UserID              uint   `gorm:"not null;unique"`
	PhoneNumber         string `gorm:"type:varchar(100)"`
	AuthenticationLevel int    `gorm:"default:0;not null"`
	BlockedLevel        int    `gorm:"default:0;not null"`
	Balance             int    `gorm:"default:0;not null"`
	IsPremium           bool   `gorm:"default:false"`
}

func (Profile) TableName() string {
	return "profiles"
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	return tx.Model(&Profile{}).Where("user_id = ?", u.ID).Update("user_id", u.ID).Error
}
