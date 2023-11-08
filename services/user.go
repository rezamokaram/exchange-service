package services

import (
	"errors"
	"fmt"
	"net/http"
	"qexchange/models"
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
	) (int, error) // returns statusCode, error

	Login(username, password string) (int, error) // returns  statusCode, error
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
	// we have access to s.db here

	// check if passwords match
	if password != passwordRepeat {
		return http.StatusBadRequest, errors.New("passwords do not match")
	}
	fmt.Println("after password match")

	// check for duplicate username or email
	var existingUser models.User
	result := s.db.Where("username = ? OR email = ?", username, email).First(&existingUser)
	if result.Error == nil {
		if existingUser.Username == username {
			return http.StatusBadRequest, errors.New("a user with this username already exists")
		}
		if existingUser.Email == email {
			return http.StatusBadRequest, errors.New("a user with this email already exists")
		}
	}
	fmt.Println("after check duplicate")

	// create user
	var newUser models.User
	newUser.Username = username
	newUser.Email = email
	fmt.Println("after check email")

	//hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	newUser.Password = string(hash)
	fmt.Println("after hash")

	//insert user into database
	if err := s.db.Create(&newUser).Error; err != nil { // Use newUser directly
		return http.StatusInternalServerError, err
	}
	fmt.Println("after insert")

	// Create a Profile for the new user
	profile := models.Profile{UserID: newUser.ID}
	if err := s.db.Create(&profile).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	fmt.Println("after create profile")

	// generate JWT token
	token, err := utils.GenerateJWTToken(newUser)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	newUser.Token = token

	fmt.Println("after token")

	// commit changes
	s.db.Save(&profile)
	s.db.Save(&newUser)

	fmt.Println("after commit")

	return http.StatusCreated, nil
}

func (s *userService) Login(username string, password string) (int, error) {
	// register logic here
	// we have access to s.db here
	return 0, nil
}
