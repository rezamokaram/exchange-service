package user

import (
	"errors"
)

type LoginRequest struct {
	Username string `json:"username" example:"newUser"`
	Password string `json:"password" example:"123456"`
}

func (lr *LoginRequest) IsValid() error {
	if lr.Password == "" || lr.Username == "" {
		return errors.New("username or password not provided")
	}

	return nil
}
