package utils

import (
	"os"
	"qexchange/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(user models.User) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"adm": user.IsAdmin,
	})

	token, err := t.SignedString([]byte(os.Getenv("JWTSECRET")))

	return token, err
}
