package internal

import (
	"errors"
	"net/http"
	"os"
	"time"

	userModels "github.com/RezaMokaram/ExchangeService/models/user"

	"github.com/golang-jwt/jwt"
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

	var newUser userModels.User
	newUser.Username = username
	newUser.Email = email

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	newUser.Password = string(hash)

	if err := s.db.Create(&newUser).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	// todo: bad code!!! where is transaction !!!!!
	profile := userModels.Profile{UserID: newUser.ID}
	if err := s.db.Create(&profile).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	// todo
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

	token, err := generateJWTToken(user)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	return http.StatusOK, token, nil
}

func generateJWTToken(user userModels.User) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"adm": user.IsAdmin,
	})

	token, err := t.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return token, err
}
