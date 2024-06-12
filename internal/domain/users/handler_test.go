package users

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (t *mockService) Activate(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (t *mockService) SignUp(ctx context.Context, u *user) error {
	args := t.Called(ctx, u)
	return args.Error(0)
}

func TestHandler_SignUp(t *testing.T) {
	type mocks struct {
		service *mockService
	}

	createMocks := func() mocks {
		return mocks{
			service: new(mockService),
		}
	}

	tests := []struct {
		name           string
		input          string
		setup          func(m *mocks) *Handler
		expectedStatus int
	}{
		{
			name:  "ValidSignUp",
			input: `{"name":"Test","email":"test@test.com","password":"pa$sw0rd"}`,
			setup: func(m *mocks) *Handler {
				m.service.On("SignUp", mock.Anything, mock.Anything).Return(nil)
				return NewHandler(m.service)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:  "InvalidInput",
			input: `{"name":"Test","email":"test@test.com","password":""}`,
			setup: func(m *mocks) *Handler {
				m.service.On("SignUp", mock.Anything, mock.Anything).Return(nil)
				return NewHandler(m.service)
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:  "DuplicateEmail",
			input: `{"name":"Test","email":"test@test.com","password":"pa$sw0rd"}`,
			setup: func(m *mocks) *Handler {
				m.service.On("SignUp", mock.Anything, mock.Anything).Return(errDuplicateEmail)
				return NewHandler(m.service)
			},
			expectedStatus: http.StatusUnprocessableEntity, // could make conflict in the future
		},
		{
			name:  "SignUpError",
			input: `{"name":"Test","email":"test@test.com","password":"pa$sw0rd"}`,
			setup: func(m *mocks) *Handler {
				m.service.On("SignUp", mock.Anything, mock.Anything).Return(errors.New("unknown error"))
				return NewHandler(m.service)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := createMocks()
			h := test.setup(&m)
			server := h.Handlers()

			request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(test.input))
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)

			assert.Equal(t, test.expectedStatus, response.Code)
		})
	}
}
