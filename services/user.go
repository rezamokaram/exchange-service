package services

import (
	"errors"
	"net/http"

	userModels "qexchange/models/user"
	"qexchange/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Register(
		username,
		password,
		passwordAgain,
		email string,
	) (int, error)

	Login(username, password string) (int, string, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

func (s *userService) Register(
	username,
	password,
	passwordRepeat,
	email string,
) (int, error) {
	// check for duplicate username or email
	var existingUser userModels.User
	result := s.db.Where("username = ? OR email = ?", username, email).First(&existingUser)
	if result.Error == nil {
		if existingUser.Username == username {
			return http.StatusBadRequest, errors.New("a user with this username already exists")
		}
		if existingUser.Email == email {
			return http.StatusBadRequest, errors.New("a user with this email already exists")
		}
	}

	// create user
	var newUser userModels.User
	newUser.Username = username
	newUser.Email = email

	//hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	newUser.Password = string(hash)

	//insert user into database
	if err := s.db.Create(&newUser).Error; err != nil { // Use newUser directly
		return http.StatusInternalServerError, err
	}

	// Create a Profile for the new user
	profile := userModels.Profile{UserID: newUser.ID}
	if err := s.db.Create(&profile).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	// commit changes
	s.db.Save(&profile)
	s.db.Save(&newUser)

	return http.StatusCreated, nil
}

func (s *userService) Login(username, password string) (int, string, error) {
	var user userModels.User

	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusUnauthorized, "", errors.New("invalid username or password")
		}
		return http.StatusInternalServerError, "", err
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return http.StatusUnauthorized, "", errors.New("invalid username or password")
	}

	token, err := utils.GenerateJWTToken(user)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	return http.StatusOK, token, nil
}
