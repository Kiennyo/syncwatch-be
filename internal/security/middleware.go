package security

import (
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	Tokens TokenVerifier
}

func Authorize(next http.HandlerFunc, requiredScopes string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		principal := ContextGetPrincipal(r)
		if principal.Sub == nil {
			authenticationRequiredResponse(w, r)
			return
		}

		scopes := *principal.Scopes

		if !strings.Contains(scopes, requiredScopes) {
			notPermittedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (a *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		// anonymous request
		if authorizationHeader == "" {
			r = contextSetPrincipal(r, &ContextValue{})
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		contextValue, err := a.Tokens.VerifyToken(token)
		if err != nil {
			invalidAuthenticationTokenResponse(w, r)
			return
		}

		r = contextSetPrincipal(r, contextValue)

		next.ServeHTTP(w, r)
	})
}
