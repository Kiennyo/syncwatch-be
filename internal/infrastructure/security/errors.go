package security

import (
	"log/slog"
	"net/http"

	"github.com/kiennyo/syncwatch-be/internal/infrastructure/json"
)

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := json.Envelope{"error": message}

	err := json.WriteJSON(w, status, env, nil)
	if err != nil {
		slog.Error(err.Error(), "request_method", r.Method, "request_url", r.URL.String())
		w.WriteHeader(500)
	}
}

func invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	ErrorResponse(w, r, http.StatusForbidden, message)
}
