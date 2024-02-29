package security

import (
	"context"
	"net/http"
)

type contextKey string

type contextValue struct {
	Sub    *string
	Scopes *string
}

const principalContext = contextKey("principal")

func contextSetPrincipal(r *http.Request, user *contextValue) *http.Request {
	ctx := context.WithValue(r.Context(), principalContext, user)
	return r.WithContext(ctx)
}

func contextGetPrincipal(r *http.Request) *contextValue {
	user, ok := r.Context().Value(principalContext).(*contextValue)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
