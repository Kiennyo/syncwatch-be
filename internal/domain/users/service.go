package users

import (
	"context"
	"log/slog"

	"github.com/kiennyo/syncwatch-be/internal/mail"
	"github.com/kiennyo/syncwatch-be/internal/security"
	"github.com/kiennyo/syncwatch-be/internal/worker"
)

type Service interface {
	SignUp(ctx context.Context, u *user) error
}

type userService struct {
	repository   Repository
	tokenCreator security.TokenCreator
	mailer       *mail.Mailer

	worker *worker.Worker
}

var _ Service = (*userService)(nil)

func NewService(r Repository, t security.TokenCreator, w *worker.Worker, m *mail.Mailer) Service {
	return &userService{
		repository:   r,
		tokenCreator: t,
		worker:       w,
		mailer:       m,
	}
}

func (s *userService) SignUp(ctx context.Context, u *user) error {
	err := s.repository.Create(ctx, u)
	if err != nil {
		return err
	}

	token, err := s.tokenCreator.CreateToken(u.ID.String(), u.Scopes, security.Activation)
	if err != nil {
		return err
	}

	s.worker.Background(func() {
		activationData := map[string]any{
			"activationToken": token,
		}

		err = s.mailer.Send(u.Email, "user_welcome.gohtml", activationData)
		if err != nil {
			slog.Error("Failed to send activation email", "reason", err.Error())
		}
	})

	return err
}
