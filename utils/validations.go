package utils

import (
	"regexp"
	"unicode"
)

// returns true if email is valid
func ValidateEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

func IsPasswordSecure(password string) bool {
	var (
		hasMinLength   = false
		hasUpper       = false
		hasLower       = false
		hasNumber      = false
		hasSpecialChar = false
		minLength      = 8 // You can adjust this as needed
	)

	if len(password) >= minLength {
		hasMinLength = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecialChar = true
		}
	}

	return hasMinLength && hasUpper && hasLower && hasNumber && hasSpecialChar
}
