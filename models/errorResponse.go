package models

type ErrorResponse struct {
	Message     string `json:"message"`
	Error      	string `json:"error"`
}

func NewErrorRespone(message string, err error) ErrorResponse {
	return ErrorResponse{
		Message: message,
		Error: err.Error(),
	}
}