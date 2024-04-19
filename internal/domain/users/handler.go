package users

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	httperr "github.com/kiennyo/syncwatch-be/internal/http/error"
	"github.com/kiennyo/syncwatch-be/internal/http/json"
	"github.com/kiennyo/syncwatch-be/internal/security"
	"github.com/kiennyo/syncwatch-be/internal/validator"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Handlers() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.signUp)
	r.Patch("/{userID}/activated", security.Authorize(h.activate, "user:activate"))

	return r
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.ReadJSON(w, r, &input)
	if err != nil {
		httperr.InvalidJSON(w, r, err)
		return
	}

	v := validator.New()

	u := &user{
		Name:  input.Name,
		Email: input.Email,
	}

	err = u.Password.set(input.Password)
	if err != nil {
		httperr.Internal(w, r, err)
		return
	}

	if validateUserInput(v, u); !v.Valid() {
		httperr.Validation(w, r, v.Errors())
		return
	}

	if err = h.service.SignUp(r.Context(), u); err != nil {
		switch {
		case errors.Is(err, errDuplicateEmail):
			// could make conflict in the future
			v.AddError("email", "a user with this email address already exists")
			httperr.Validation(w, r, v.Errors())
		default:
			httperr.Internal(w, r, err)
		}
		return
	}

	err = json.WriteJSON(w, http.StatusCreated, json.Envelope{"user": u}, nil)
	if err != nil {
		httperr.Internal(w, r, err)
	}
}

func (h *Handler) activate(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	principal := security.ContextGetPrincipal(r)

	if userID != principal.Sub {
		httperr.Forbidden(w, r)
	}

	err := h.service.Activate(r.Context(), principal.Sub)
	if err != nil {
		httperr.Internal(w, r, err)
		return
	}

	err = json.WriteJSON(w, http.StatusOK, json.Envelope{"message": "User activated successfully"}, nil)
	if err != nil {
		httperr.Internal(w, r, err)
	}
}
