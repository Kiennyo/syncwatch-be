package auth

import (
	"github.com/go-chi/chi/v5"
)

func Handlers() chi.Router {
	r := chi.NewRouter()
	return r
}
