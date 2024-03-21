package user

type RegisterRequest struct {
	Username       string `json:"username" example:"newUser"`
	Email          string `json:"email" example:"newUser@example.com"`
	Password       string `json:"password" example:"123456"`
	PasswordRepeat string `json:"password_repeat" example:"123456"`
}
