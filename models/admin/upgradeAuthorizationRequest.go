package admin

type UpdateAuthenticationLevelRequest struct {
	Username     string `json:"username" example:"user2"`
	NewAuthLevel int    `json:"new_auth_level" example:"0"`
}