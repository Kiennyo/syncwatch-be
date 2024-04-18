package security

import (
	"net/http"

	"github.com/kiennyo/syncwatch-be/internal/http/error"
)

func invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	error.Response(w, r, http.StatusUnauthorized, message)
}

func authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	error.Response(w, r, http.StatusUnauthorized, message)
}

func notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	error.Response(w, r, http.StatusForbidden, message)
}
