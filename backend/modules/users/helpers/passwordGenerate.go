package helpers

import (
	"crypto/rand"
	"math/big"
	"strings"

	"elogika.vsb.cz/backend/models"
)

const (
	passwordLength = 15
	upperLetters   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerLetters   = "abcdefghijklmnopqrstuvwxyz"
	numbers        = "0123456789"
	specialChars   = "!@#$%^&*()-_=+[]{}|;:'\",.<>?/`~"
	allChars       = upperLetters + lowerLetters + numbers + specialChars
)

func GeneratePassword(user *models.User) string {
	for {
		var password strings.Builder
		// Ensure at least one of each required type
		password.WriteByte(upperLetters[randomIndex(len(upperLetters))])
		password.WriteByte(numbers[randomIndex(len(numbers))])
		password.WriteByte(specialChars[randomIndex(len(specialChars))])

		// Fill the rest with random characters
		for password.Len() < passwordLength {
			password.WriteByte(allChars[randomIndex(len(allChars))])
		}

		// Shuffle password to avoid predictable pattern
		pwRunes := []rune(password.String())
		shuffleRunes(pwRunes)
		finalPassword := string(pwRunes)

		// Check if it contains name or surname
		lowerPassword := strings.ToLower(finalPassword)
		if user.FirstName != "" && strings.Contains(lowerPassword, strings.ToLower(user.FirstName)) {
			continue
		}
		if user.FamilyName != "" && strings.Contains(lowerPassword, strings.ToLower(user.FamilyName)) {
			continue
		}

		return finalPassword
	}
}

// randomIndex returns a secure random index from 0 to max-1
func randomIndex(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

// shuffleRunes shuffles a slice of runes in-place
func shuffleRunes(runes []rune) {
	for i := len(runes) - 1; i > 0; i-- {
		j := randomIndex(i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}
}
