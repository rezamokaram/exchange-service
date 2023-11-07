package services

import (
	"qexchange/models"

	"gorm.io/gorm"
)

type UserService interface {
	Register(username, password, passwordAgain string) (*models.User, error)
	Login(username, password string) (*models.User, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

func (s *userService) Register(username, password, passwordAgain string) (*models.User, error) {
	// register logic here
	// we have access to s.db here

	return nil, nil
}

func (s *userService) Login(username string, password string) (*models.User, error) {
	// register logic here
	// we have access to u.db here
	return nil, nil
}
