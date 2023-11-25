package services

import (
	"errors"
	"net/http"
	"os"
	"qexchange/models"

	"gorm.io/gorm"
)

type AdminService interface {
	UpgradeToAdmin(user models.User, adminPasswordJSON string) error
	UpdateAuthenticationLevel(username string, newAuthLevel int) error
	BlockUser(username string, temporary bool) (int, error)
	UnblockUser(username string) (int, error)
	GetUserInfo(username string) (models.UserInfo, int, error)
}

type adminService struct {
	db *gorm.DB
}

const (
	Authenticated = iota
	Unauthenticated
)

const (
	Unblocked = iota
	BlockedTemporarily
	BlockedPermanently
)

func NewAdminService(db *gorm.DB) AdminService {
	return &adminService{
		db: db,
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

func (s *adminService) BlockUser(username string, temporary bool) (int, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).Preload("Profile").First(&user).Error; err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	if temporary && user.Profile.BlockedLevel == BlockedTemporarily {
		return http.StatusBadRequest, errors.New("user id already temporarily blocked")
	}

	if !temporary && user.Profile.BlockedLevel == BlockedPermanently {
		return http.StatusBadRequest, errors.New("user id already permanently blocked")
	}

	var newBlockedLevel int
	if temporary {
		newBlockedLevel = BlockedTemporarily
	} else {
		newBlockedLevel = BlockedPermanently
	}

	user.Profile.BlockedLevel = newBlockedLevel

	if err := s.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user).Error; err != nil {
		return http.StatusBadRequest, errors.New("failed updating user")
	}

	// instead of the above we can also do this do update the profile association
	// err := s.db.Model(&models.Profile{}).Where("user_id = ?", user.ID).Update("blocked_level", newBlockedLevel).Error

	return http.StatusOK, nil
}

func (s *adminService) UnblockUser(username string) (int, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).Preload("Profile").First(&user).Error; err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	if user.Profile.BlockedLevel == Unblocked {
		return http.StatusBadRequest, errors.New("user is not blocked")
	}

	user.Profile.BlockedLevel = Unblocked

	if err := s.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user).Error; err != nil {
		return http.StatusBadRequest, errors.New("failed updating user")
	}

	// instead of the above we can also do this do update the profile association
	// err := s.db.Model(&models.Profile{}).Where("user_id = ?", user.ID).Update("blocked_level", Unblocked).Error

	return http.StatusOK, nil
}

func (s *adminService) GetUserInfo(username string) (models.UserInfo, int, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).Preload("Profile").Preload("BankingInfo").First(&user).Error; err != nil {
		return models.UserInfo{}, http.StatusNotFound, errors.New("user not found")
	}

	newUserInfo := models.NewUserInfo(user)

	return newUserInfo, http.StatusOK, nil
}
