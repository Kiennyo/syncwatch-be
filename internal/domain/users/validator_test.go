package users

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiennyo/syncwatch-be/internal/validator"
)

func TestIsValidPasswordComposition(t *testing.T) {
	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{
			name:     "all three types present",
			password: "Pass13!",
			valid:    true,
		},
		{
			name:     "missing number",
			password: "PassOnly!",
			valid:    false,
		},
		{
			name:     "missing special char",
			password: "Pass123",
			valid:    false,
		},
		{
			name:     "missing letter",
			password: "123!@#",
			valid:    false,
		},
		{
			name:     "only number",
			password: "123456",
			valid:    false,
		},
		{
			name:     "only special char",
			password: "!@#$$%",
			valid:    false,
		},
		{
			name:     "empty string",
			password: "",
			valid:    false,
		},
		{
			name:     "single character, a letter",
			password: "a",
			valid:    false,
		},
		{
			name:     "single character, a number",
			password: "1",
			valid:    false,
		},
		{
			name:     "single character, a special character",
			password: "!",
			valid:    false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := isValidPasswordComposition(tc.password)
			assert.Equal(t, tc.valid, actual, "unexpected result for password: %s", tc.password)
		})
	}
}

func TestValidateUserInput(t *testing.T) {
	cases := []struct {
		name             string
		user             *user
		password         string
		validationErrors map[string]string
	}{
		{
			name:     "Empty",
			user:     &user{},
			password: "",
			validationErrors: map[string]string{
				"name":     "must be provided",
				"email":    "must be provided",
				"password": "must be provided",
			},
		},
		{
			name:     "ExceedingLength",
			user:     &user{Name: strings.Repeat("a", 501), Email: "user@example.com"},
			password: "pa$$w0rd",
			validationErrors: map[string]string{
				"name": "must not be more than 500 bytes long",
			},
		},
		{
			name:     "InvalidEmailFormat",
			user:     &user{Name: "User", Email: "user@example"},
			password: "pa$$w0rd",
			validationErrors: map[string]string{
				"email": "must be a valid email address",
			},
		},
		{
			name:     "ShortPassword",
			user:     &user{Name: "User", Email: "user@example.com"},
			password: "1",
			validationErrors: map[string]string{
				"password": "must be at least 8 bytes long",
			},
		},
		{
			name:     "TooLongPassword",
			user:     &user{Name: "User", Email: "user@example.com"},
			password: strings.Repeat("pa$$w0rd00", 10),
			validationErrors: map[string]string{
				"password": "must not be more than 72 bytes long",
			},
		},
		{
			name:             "ValidInput",
			user:             &user{Name: "User", Email: "user@example.com"},
			password:         "pa$$w0rd",
			validationErrors: map[string]string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := validator.New()
			tc.user.Password = password{plaintext: &tc.password} //nolint
			validateUserInput(v, tc.user)
			assert.Equal(t, tc.validationErrors, v.Errors())
		})
	}
}
