package users

import (
	"regexp"
	"unicode"

	"github.com/kiennyo/syncwatch-be/internal/validator"
)

// nolint
var emailRX = regexp.MustCompile(".+@.+\\..+")

func validateUserInput(v *validator.Validator, u *user) {
	v.Check(u.Name != "", "name", "must be provided")
	v.Check(len(u.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(u.Email != "", "email", "must be provided")
	v.Check(validator.Matches(u.Email, emailRX), "email", "must be a valid email address")

	v.Check(*u.Password.plaintext != "", "password", "must be provided")
	v.Check(len(*u.Password.plaintext) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(*u.Password.plaintext) <= 72, "password", "must not be more than 72 bytes long")
	v.Check(isValidPasswordComposition(*u.Password.plaintext), "password", "must be a valid password")
}

func isValidPasswordComposition(password string) bool {
	hasNumber, hasSpecialChar, hasLetter := false, false, false

	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecialChar = true
		case unicode.IsLetter(c) || c == ' ':
			hasLetter = true
		default:
			return false
		}
	}

	return hasNumber && hasSpecialChar && hasLetter
}
