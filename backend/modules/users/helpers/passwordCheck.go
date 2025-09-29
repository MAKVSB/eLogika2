package helpers

import (
	"fmt"
	"strings"
	"unicode"

	"elogika.vsb.cz/backend/models"
)

func PasswordCheck(password string, user *models.User) (bool, string) {
	const minLength = 8
	specialChars := "!@#$%^&*()-_=+[]{}|;:'\",.<>?/`~"

	if len(password) < minLength {
		return false, fmt.Sprintf("Password must be at least %d characters long", minLength)
	}

	var hasNumber, hasUpper, hasSpecial bool
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpper = true
		case strings.ContainsRune(specialChars, c):
			hasSpecial = true
		}
	}

	if !hasNumber {
		return false, "Password must contain at least one number"
	}
	if !hasUpper {
		return false, "Password must contain at least one uppercase letter"
	}
	if !hasSpecial {
		return false, "Password must contain at least one special character"
	}

	lowerPassword := strings.ToLower(password)
	if user.FirstName != "" && strings.Contains(lowerPassword, strings.ToLower(user.FirstName)) {
		return false, "Password cannot contain your name"
	}
	if user.FamilyName != "" && strings.Contains(lowerPassword, strings.ToLower(user.FamilyName)) {
		return false, "Password cannot contain your surname"
	}

	return true, ""
}
