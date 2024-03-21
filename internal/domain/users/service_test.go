package users

//import (
//	"context"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//
//	"github.com/kiennyo/syncwatch-be/internal/domain/users/mocks"
//
//	"github.com/kiennyo/syncwatch-be/internal/worker"
//)
//
//func TestUserService_SignUp(t *testing.T) {
//	ctx := context.Background()
//
//	t.Parallel()
//
//	tt := []struct {
//		name    string
//		setup   func() *userService
//		user    *user
//		wantErr bool
//	}{
//		{
//			name: "Success",
//			setup: func() *userService {
//				repo := mocks.NewMockRepository()
//				tokenCreator := newMockTokenCreator()
//				mailer := newMockMailer()
//				wrk := worker.New()
//
//				return &userService{
//					repository:   repo,
//					tokenCreator: tokenCreator,
//					mailer:       mailer,
//					worker:       wrk,
//				}
//			},
//			user:    &user{},
//			wantErr: false,
//		},
//{
//	name: "RepoError",
//	setup: func() *userService {
//		repo := mocks.NewMockRepository()
//		repo.
//			repo.On("Create", ctx, &user{}).Return(errors.New("some error"))
//		tokenCreator := newMockTokenCreator()
//		mailer := newMockMailer()
//		wrk := newMockWorker()
//
//		return &userService{
//			repository:   repo,
//			tokenCreator: tokenCreator,
//			mailer:       mailer,
//			worker:       wrk,
//		}
//	},
//	user:    &user{},
//	wantErr: true,
//},
//{
//	name: "TokenCreatorError",
//	setup: func() *userService {
//		repo := newMockRepo()
//		tokenCreator := newMockTokenCreator()
//		tokenCreator.On("CreateToken", "", []string{}, security.Activation).Return("", errors.New("some error"))
//		mailer := newMockMailer()
//		wrk := newMockWorker()
//
//		return &userService{
//			repository:   repo,
//			tokenCreator: tokenCreator,
//			mailer:       mailer,
//			worker:       wrk,
//		}
//	},
//	user:    &user{},
//	wantErr: true,
//},
//{
//	name: "MailerError",
//	setup: func() *userService {
//		repo := newMockRepo()
//		tokenCreator := newMockTokenCreator()
//		mailer := newMockMailer()
//		mailer.On("Send", "", "user_welcome.gohtml", map[string]any{
//			"activationToken": "",
//		}).Return(errors.New("some error"))
//		wrk := newMockWorker()
//
//		return &userService{
//			repository:   repo,
//			tokenCreator: tokenCreator,
//			mailer:       mailer,
//			worker:       wrk,
//		}
//	},
//	user:    &user{},
//	wantErr: true,
//},
//	}
//
//	for _, tc := range tt {
//		t.Run(tc.name, func(t *testing.T) {
//			sut := tc.setup()
//
//			err := sut.SignUp(ctx, tc.user)
//
//			if tc.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
