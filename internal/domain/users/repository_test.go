package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiennyo/syncwatch-be/internal/testhelpers"
)

func TestUserRepository_Create(t *testing.T) {
	ctx := context.Background()
	container, err := testhelpers.CreateTestDB(ctx)

	assert.NotNil(t, container)
	assert.Nil(t, err)
	assert.NotNil(t, container.DB)

	repository := NewRepository(container.DB)
	u := &user{
		Name:   "John",
		Email:  "test@test.lt",
		Scopes: nil,
	}
	err = u.Password.set("test")
	assert.Nil(t, err)

	err = repository.Create(ctx, u)
	assert.Nil(t, err)

	assert.Equal(t, []string{"user:activate"}, u.Scopes)
}

// this could also be in a single test
func TestUserRepository_CreateUserWithExistingEmail(t *testing.T) {
	ctx := context.Background()
	container, err := testhelpers.CreateTestDB(ctx)

	assert.NotNil(t, container)
	assert.Nil(t, err)
	assert.NotNil(t, container.DB)

	repository := NewRepository(container.DB)
	u := &user{
		Name:   "John",
		Email:  "test@test.com",
		Scopes: nil,
	}
	err = u.Password.set("test")
	assert.Nil(t, err)

	err = repository.Create(ctx, u)
	assert.Nil(t, err)
	assert.Equal(t, []string{"user:activate"}, u.Scopes)

	err = repository.Create(ctx, u)
	assert.Equal(t, errDuplicateEmail, err)
}
