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
	Activate(ctx context.Context, id string) error
}

type userService struct {
	repository   Repository
	tokenCreator security.TokenCreator
	mailer       mail.Sender
}

var _ Service = (*userService)(nil)

func NewService(r Repository, t security.TokenCreator, m mail.Sender) Service {
	return &userService{
		repository:   r,
		tokenCreator: t,
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

	worker.Background(func() {
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

func (s *userService) Activate(ctx context.Context, id string) error {
	usr, err := s.repository.FindById(ctx, id)
	if err != nil {
		return err
	}

	usr.Activated = true

	err = s.repository.Activate(ctx, usr)
	if err != nil {
		return err
	}

	return nil
}
