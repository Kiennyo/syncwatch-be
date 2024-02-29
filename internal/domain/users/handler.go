package users

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kiennyo/syncwatch-be/internal/infrastructure/security"
)

type Handler struct {
	// goal just to pass service and call it a day
}

func (h *Handler) Handlers() chi.Router {
	r := chi.NewRouter()
	r.Get("/", security.Authorize(hello, "OMFG"))
	return r
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("111"))
}

func NewHandler() Handler {
	return Handler{}
}
