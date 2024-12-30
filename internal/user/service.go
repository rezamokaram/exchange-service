package user

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/RezaMokaram/ExchangeService/internal/user/domain"
	"github.com/RezaMokaram/ExchangeService/internal/user/port"
)

var (
	ErrUserOnCreate           = errors.New("error on creating new user")
	ErrUserCreationValidation = errors.New("validation failed")
	ErrUserNotFound           = errors.New("user not found")
)

type service struct {
	repo port.Repo
}

func NewService(repo port.Repo) port.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateUser(ctx context.Context, user domain.User) (domain.UserID, error) {
	if err := user.Validate(); err != nil {
		return 0, fmt.Errorf("%w %w", ErrUserCreationValidation, err)
	}

	userID, err := s.repo.Create(ctx, user)
	if err != nil {
		log.Println("error on creating new user : ", err.Error())
		return 0, ErrUserOnCreate
	}

	return userID, nil
}

func (s *service) GetUserByFilter(ctx context.Context, filter *domain.UserFilter) (*domain.User, error) {
	user, err := s.repo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// func (s *service) Register(
// 	username,
// 	password,
// 	passwordRepeat,
// 	email string,
// ) (int, error) {
// 	var existingUser userModels.User
// 	result := s.db.Where("username = ? OR email = ?", username, email).First(&existingUser)
// 	if result.Error == nil {
// 		if existingUser.Username == username {
// 			return http.StatusBadRequest, errors.New("a user with this username already exists")
// 		}
// 		if existingUser.Email == email {
// 			return http.StatusBadRequest, errors.New("a user with this email already exists")
// 		}
// 	}

// 	var newUser userModels.User
// 	newUser.Username = username
// 	newUser.Email = email

// 	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return http.StatusInternalServerError, err
// 	}
// 	newUser.Password = string(hash)

// 	if err := s.db.Create(&newUser).Error; err != nil {
// 		return http.StatusInternalServerError, err
// 	}

// 	// todo: bad code!!! where is transaction !!!!!
// 	profile := userModels.Profile{UserID: newUser.ID}
// 	if err := s.db.Create(&profile).Error; err != nil {
// 		return http.StatusInternalServerError, err
// 	}

// 	// todo
// 	s.db.Save(&profile)
// 	s.db.Save(&newUser)

// 	return http.StatusCreated, nil
// }

// func (s *service) Login(username, password string) (int, string, error) {
// 	var user userModels.User

// 	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return http.StatusUnauthorized, "", errors.New("invalid username or password")
// 		}
// 		return http.StatusInternalServerError, "", err
// 	}

// 	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
// 	if err != nil {
// 		return http.StatusUnauthorized, "", errors.New("invalid username or password")
// 	}

// 	token, err := generateJWTToken(user)
// 	if err != nil {
// 		return http.StatusInternalServerError, "", err
// 	}

// 	return http.StatusOK, token, nil
// }

// func generateJWTToken(user userModels.User) (string, error) {
// 	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"id":  user.ID,
// 		"exp": time.Now().Add(time.Hour * 72).Unix(),
// 		"adm": user.IsAdmin,
// 	})

// 	token, err := t.SignedString([]byte(os.Getenv("JWT_SECRET")))

// 	return token, err
// }
