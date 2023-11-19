package services

import (
	"errors"
	"gorm.io/gorm"
	"os"
	"qexchange/database"
	"qexchange/models"
)

type AdminService interface {
	UpgradeToAdmin(user models.User, adminPasswordJSON string) error
	UpdateAuthenticationLevel(username string, newAuthLevel int) error
}

type adminService struct {
	db        *gorm.DB
	dbService database.DataBaseService
}

const (
	Unauthenticated = iota
	Authenticated
)

func NewAdminService(db *gorm.DB) AdminService {
	return &adminService{
		db:        db,
		dbService: database.NewDBService(db),
	}
}

func (s *adminService) UpgradeToAdmin(user models.User, adminPasswordJSON string) error {

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPasswordJSON != adminPassword {
		return errors.New("invalid admin password")
	}

	user.IsAdmin = true
	return s.db.Save(&user).Error
}

func (s *adminService) UpdateAuthenticationLevel(username string, newAuthLevel int) error {
	if newAuthLevel != Unauthenticated && newAuthLevel != Authenticated {
		return errors.New("newAuthLevel is invalid")
	}

	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}

	if user.IsAdmin {
		return errors.New("user is already admin")
	}

	err := s.db.Model(&models.Profile{}).Where("user_id = ?", user.ID).Update("authentication_level", newAuthLevel).Error
	return err
}
