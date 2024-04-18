package security

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kiennyo/syncwatch-be/internal/config"
)

func TestTokensFactory_CreateToken(t *testing.T) {
	factory := NewTokenFactory(config.Security{
		JWTSecret: "mock_secret",
		Iss:       "syncwatch.io",
		Aud:       "syncwatch.io",
	})

	// Test cases
	tests := []struct {
		name   string
		userID string
		scopes []string
		exp    Expiration
	}{
		{name: "3 days", userID: uuid.New().String(), scopes: []string{"user:activate"}, exp: Activation},
		{name: "15 mins", userID: uuid.New().String(), scopes: []string{"user:view", "user:edit"}, exp: Access},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := factory.CreateToken(tt.userID, tt.scopes, tt.exp)

			// Check error and token
			assert.Nil(t, err)
			assert.NotEmpty(t, tokenString)

			// Verify created token
			contextValue, err := factory.VerifyToken(tokenString)

			// Check error and returned claims
			assert.Nil(t, err)
			assert.Equal(t, tt.userID, *contextValue.Sub)
			assert.Equal(t, strings.Join(tt.scopes, " "), *contextValue.Scopes)
		})
	}
}
