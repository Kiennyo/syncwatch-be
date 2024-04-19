package security

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/kiennyo/syncwatch-be/internal/config"
)

type Expiration time.Duration

const (
	Activation = Expiration(3 * 24 * time.Hour)
	Access     = Expiration(15 * time.Minute)
)

type TokenCreator interface {
	CreateToken(userID string, scopes []string, exp Expiration) (string, error)
}

type TokenVerifier interface {
	VerifyToken(token string) (*ContextValue, error)
}

type TokensFactory struct {
	secret []byte
	iss    string
	aud    string
}

type Claims struct {
	Scopes string `json:"scopes"`
	jwt.RegisteredClaims
}

var (
	_ TokenCreator  = (*TokensFactory)(nil)
	_ TokenVerifier = (*TokensFactory)(nil)
)

func NewTokenFactory(cfg config.Security) *TokensFactory {
	return &TokensFactory{
		secret: []byte(cfg.JWTSecret),
		iss:    cfg.Iss,
		aud:    cfg.Aud,
	}
}

func (t *TokensFactory) CreateToken(userID string, scopes []string, exp Expiration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Audience:  []string{t.aud},
				Subject:   userID,
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(exp))),
				ID:        uuid.New().String(),
				Issuer:    t.iss,
			},
			Scopes: strings.Join(scopes, " "),
		})

	tokenString, err := token.SignedString(t.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *TokensFactory) VerifyToken(token string) (*ContextValue, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, func(_ *jwt.Token) (any, error) {
		return t.secret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithAudience(t.aud),
		jwt.WithIssuer(t.iss),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*Claims)
	if !ok {
		return nil, err
	}

	return &ContextValue{
		Sub:    claims.Subject,
		Scopes: claims.Scopes,
	}, nil
}
