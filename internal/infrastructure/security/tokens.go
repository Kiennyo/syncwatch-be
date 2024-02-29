package security

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Tokens struct {
	secret []byte
}

type Claims struct {
	Scopes string `json:"scopes"`
	jwt.RegisteredClaims
}

func (t *Tokens) CreateToken(userID int64, scopes string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Audience:  []string{"syncwatch.io"},
				Subject:   strconv.FormatInt(userID, 10),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
				ID:        uuid.New().String(),
				Issuer:    "syncwatch.io",
			},
			Scopes: scopes,
		})

	tokenString, err := token.SignedString(t.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *Tokens) VerifyToken(token string) (*string, *string, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return t.secret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithAudience("syncwatch.io"),
		jwt.WithIssuer("syncwatch.io"),
	)

	switch {
	case jwtToken.Valid:
		fmt.Println("You look nice today")
	case errors.Is(err, jwt.ErrTokenMalformed):
		fmt.Println("That's not even a token")
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		fmt.Println("Invalid signature")
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		fmt.Println("Timing is everything")
	default:
		fmt.Println("Couldn't handle this token:", err)
	}

	if err != nil {
		return nil, nil, err
	}

	if !jwtToken.Valid {
		return nil, nil, fmt.Errorf("invalid token")
	}

	claims, ok := jwtToken.Claims.(*Claims)
	if !ok {
		return nil, nil, err
	}

	return &claims.Subject, &claims.Scopes, nil
}

func NewTokens(secret string) *Tokens {
	return &Tokens{
		secret: []byte(secret),
	}
}
