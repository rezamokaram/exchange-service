package user

type LoginRequest struct {
	Username string `json:"username" example:"newUser"`
	Password string `json:"password" example:"123456"`
}
