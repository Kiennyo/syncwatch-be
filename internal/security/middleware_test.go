package security

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kiennyo/syncwatch-be/internal/config"
)

//nolint:revive,cognitive-complexity
func TestAuthMiddleware_Authenticate(t *testing.T) {
	tokenFactory := NewTokenFactory(config.Security{
		JWTSecret: "superSecret",
		Iss:       "syncwatch.io",
		Aud:       "syncwatch.io",
	})

	subject := uuid.New().String()
	scopes := []string{"users:activate", "users:view"}
	scopesJoined := strings.Join(scopes, " ")

	token, err := tokenFactory.CreateToken(subject, scopes, Activation)
	assert.Nil(t, err)

	testCases := []struct {
		name               string
		authorizationToken string
		expectedStatusCode int
		context            *ContextValue
	}{
		{
			name:               "Anonymous request",
			authorizationToken: "",
			expectedStatusCode: http.StatusOK,
			context:            &ContextValue{},
		},
		{
			name:               "Invalid Token",
			authorizationToken: "Bearer Random.Invalid.Token",
			expectedStatusCode: http.StatusUnauthorized,
			context:            nil,
		},
		{
			name:               "Invalid Schema",
			authorizationToken: "Basic Random.Invalid.Token",
			expectedStatusCode: http.StatusUnauthorized,
			context:            nil,
		},
		{
			name:               "Valid Token",
			authorizationToken: fmt.Sprintf("Bearer %s", token),
			expectedStatusCode: http.StatusOK,
			context: &ContextValue{
				Sub:    subject,
				Scopes: scopesJoined,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			res := httptest.NewRecorder()

			if tc.authorizationToken != "" {
				req.Header.Set("Authorization", tc.authorizationToken)
			}

			am := &AuthMiddleware{
				Tokens: tokenFactory,
			}

			nextHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				if tc.context != nil {
					principal := ContextGetPrincipal(r)
					assert.Equal(t, tc.context.Scopes, principal.Scopes)
					assert.Equal(t, tc.context.Sub, principal.Sub)
				}
			})

			am.Authenticate(nextHandler).ServeHTTP(res, req)
			assert.Equal(t, tc.expectedStatusCode, res.Result().StatusCode) // nolint
		})
	}
}

func TestAuthMiddleware_Authorize(t *testing.T) {
	sub := uuid.New().String()
	tests := []struct {
		name            string
		requiredScopes  string
		principalScopes string
		subject         string
		expectedStatus  int
	}{
		{
			name:            "Request scopes matches with handler's scopes",
			requiredScopes:  "user:read",
			principalScopes: "user:read",
			subject:         sub,
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "Request scopes doesnt match with handler's scopes",
			requiredScopes:  "admin:write",
			principalScopes: "admin:read",
			subject:         sub,
			expectedStatus:  http.StatusForbidden,
		},
		{
			name:            "Request doesn't have required scopes",
			requiredScopes:  "admin:read",
			principalScopes: "",
			subject:         sub,
			expectedStatus:  http.StatusForbidden,
		},
		{
			name:            "Request doesn't have subject",
			requiredScopes:  "admin:read",
			principalScopes: "",
			subject:         "",
			expectedStatus:  http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//nolint:gosec,G601
			ctx := context.WithValue(context.TODO(), principalContext, &ContextValue{
				Sub:    test.subject,
				Scopes: test.principalScopes,
			})
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
			Authorize(handler, test.requiredScopes).ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatus, rr.Code)
		})
	}
}
