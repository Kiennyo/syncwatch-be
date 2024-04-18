package users

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kiennyo/syncwatch-be/internal/security"
	"github.com/kiennyo/syncwatch-be/internal/worker"
)

// mocks
type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) Create(ctx context.Context, u *user) error {
	args := r.Called(ctx, u)
	return args.Error(0)
}

type tokenCreatorMock struct {
	mock.Mock
}

func (t *tokenCreatorMock) CreateToken(userID string, scopes []string, exp security.Expiration) (string, error) {
	args := t.Called(userID, scopes, exp)
	return args.String(0), args.Error(1)
}

type mailSenderMock struct {
	mock.Mock
}

func (m *mailSenderMock) Send(recipient, templateFile string, data any) error {
	args := m.Called(recipient, templateFile, data)
	return args.Error(0)
}

//nolint:revive,function-length
func TestUserService_SignUp(t *testing.T) {
	ctx := context.Background()

	type mocks struct {
		repo         *repositoryMock
		tokenCreator *tokenCreatorMock
		mailSender   *mailSenderMock
	}

	createMocks := func() mocks {
		return mocks{
			repo:         new(repositoryMock),
			tokenCreator: new(tokenCreatorMock),
			mailSender:   new(mailSenderMock),
		}
	}

	tt := []struct {
		name    string
		setup   func(m *mocks) Service
		user    *user
		wantErr bool
	}{
		{
			name: "Success",
			setup: func(m *mocks) Service {
				m.repo.On("Create", mock.Anything, mock.Anything).Return(nil)
				m.tokenCreator.On("CreateToken", uuid.Nil.String(), []string(nil), security.Activation).
					Return("token", nil)
				m.mailSender.On("Send", "email@test.com", "user_welcome.gohtml", map[string]any{
					"activationToken": "token",
				}).Return(nil)

				return NewService(m.repo, m.tokenCreator, m.mailSender)
			},
			user: &user{
				Email: "email@test.com",
			},
			wantErr: false,
		},
		{
			name: "RepoError",
			setup: func(m *mocks) Service {
				m.repo.On("Create", ctx, &user{}).Return(errors.New("some error"))

				return NewService(m.repo, m.tokenCreator, m.mailSender)
			},
			user:    &user{},
			wantErr: true,
		},
		{
			name: "TokenCreatorError",
			setup: func(m *mocks) Service {
				m.repo.On("Create", ctx, mock.Anything).Return(nil)
				m.tokenCreator.On("CreateToken", uuid.Nil.String(), []string{"scope"}, security.Activation).
					Return("", errors.New("some error"))

				return NewService(m.repo, m.tokenCreator, m.mailSender)
			},
			user: &user{
				ID:     uuid.Nil,
				Scopes: []string{"scope"},
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := createMocks()
			sut := tc.setup(&m)

			err := sut.SignUp(ctx, tc.user)

			worker.Wait()

			m.mailSender.AssertExpectations(t)
			m.repo.AssertExpectations(t)
			m.tokenCreator.AssertExpectations(t)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
