package admin


type BlockUserRequest struct {
	Username  string `json:"username" example:"user1"`
	Temporary bool   `json:"temporary" example:"false"`
}