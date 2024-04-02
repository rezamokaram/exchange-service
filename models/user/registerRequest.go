package user

import (
	"errors"
)

type RegisterRequest struct {
	Username       string `json:"username" example:"newUser"`
	Email          string `json:"email" example:"newUser@example.com"`
	Password       string `json:"password" example:"123456"`
	PasswordRepeat string `json:"password_repeat" example:"123456"`
}

func (rr *RegisterRequest) IsValid() error {
	if rr.Email == "" {
		return errors.New("email not provided")
	}

	if !isEmailValid(rr.Email) {
		return errors.New("email is not valid")
	}

	if rr.Username == "" {
		return errors.New("username not provided")
	}

	if rr.Password == "" || rr.PasswordRepeat == "" {
		return errors.New("password not provided")
	}

	if rr.Password != rr.PasswordRepeat {
		return errors.New("passwords do not match")
	}

	if !isPasswordSecure(rr.Password) {
		return errors.New("password is not secure")
	}

	return nil
}
