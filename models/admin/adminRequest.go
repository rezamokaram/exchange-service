package admin

type UpdateUserToAdminRequest struct {
	AdminPassword string `json:"admin_password" example:"secret"`
}