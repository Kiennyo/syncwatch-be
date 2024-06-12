package error

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/kiennyo/syncwatch-be/internal/http/json"
)

func Response(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := json.Envelope{"error": message}

	err := json.WriteJSON(w, status, env, nil)
	if err != nil {
		slog.Error(err.Error(), "request_method", r.Method, "request_url", r.URL.String())
		w.WriteHeader(500)
	}
}

func Internal(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error(err.Error(), "request_method", r.Method, "request_url", r.URL.String())

	message := "the server encountered a problem and could not process your request"
	Response(w, r, http.StatusInternalServerError, message)
}

func Validation(w http.ResponseWriter, r *http.Request, errs map[string]string) {
	Response(w, r, http.StatusUnprocessableEntity, errs)
}

func Forbidden(w http.ResponseWriter, r *http.Request) {
	Response(w, r, http.StatusForbidden, "Forbidden")
}

func InvalidJSON(w http.ResponseWriter, r *http.Request, err error) {
	var mr *json.MalformedRequest
	if errors.As(err, &mr) {
		Response(w, r, mr.Status, mr.Msg)
	} else {
		Internal(w, r, err)
	}
}
