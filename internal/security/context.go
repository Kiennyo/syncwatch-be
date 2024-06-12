package security

import (
	"context"
	"net/http"
)

type contextKey string

type ContextValue struct {
	Sub    string
	Scopes string
}

const principalContext = contextKey("principal")

func contextSetPrincipal(r *http.Request, principal *ContextValue) *http.Request {
	ctx := context.WithValue(r.Context(), principalContext, principal)
	return r.WithContext(ctx)
}

func ContextGetPrincipal(r *http.Request) *ContextValue {
	principal, ok := r.Context().Value(principalContext).(*ContextValue)
	if !ok {
		panic("missing principal value in request context")
	}
	return principal
}
