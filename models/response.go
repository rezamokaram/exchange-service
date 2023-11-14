package models

type Response struct {
	Message     string `json:"message"`
}

func NewRespone(message string) ErrorResponse {
	return ErrorResponse{
		Message: message,
	}
}