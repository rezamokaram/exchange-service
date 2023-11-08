package services

import (
	"errors"
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

	Login(username, password string) (int, string, error) // returns  statusCode, token ,error
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

	// create user
	var newUser models.User
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
	profile := models.Profile{UserID: newUser.ID}
	if err := s.db.Create(&profile).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	// generate JWT token
	token, err := utils.GenerateJWTToken(newUser)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	newUser.Token = token

	// commit changes
	s.db.Save(&profile)
	s.db.Save(&newUser)

	return http.StatusCreated, nil
}

func (s *userService) Login(username, password string) (int, string, error) {
	var user models.User

	// check for the existence of the user with the given username
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// user not found
			return http.StatusUnauthorized, "", errors.New("invalid username or password")
		}
		// other possible errors
		return http.StatusInternalServerError, "", err
	}

	// user found, now compare the given password with the hashed password in the database
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// incorrect password
		return http.StatusUnauthorized, "", errors.New("invalid username or password")
	}

	// password is correct, now generate JWT token
	token, err := utils.GenerateJWTToken(user)
	if err != nil {
		// error generating JWT token
		return http.StatusInternalServerError, "", err
	}

	// return the token and a status of OK
	return http.StatusOK, token, nil
}
