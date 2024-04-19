package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const userInactiveRole = "user-inactive"
const userActiveRole = "user-active"

type Repository interface {
	Create(ctx context.Context, u *user) error
	FindById(ctx context.Context, id string) (*user, error)
	Activate(ctx context.Context, usr *user) error
}

type userRepository struct {
	DB *pgxpool.Pool
}

var _ Repository = (*userRepository)(nil)

func NewRepository(db *pgxpool.Pool) Repository {
	return &userRepository{DB: db}
}

func (r *userRepository) Create(ctx context.Context, u *user) error {
	query := `
		WITH user_insert AS (
			INSERT INTO "user" (name, email, password_hash, updated_at, role_id)
				VALUES (@name, @email, @password_hash, NOW(),
						(SELECT id FROM role WHERE slug = @role))
				RETURNING id, created_at, role_id, updated_at)
		SELECT user_insert.id, user_insert.created_at, user_insert.updated_at, JSON_AGG(permission.slug)
		FROM user_insert
		INNER JOIN role ON user_insert.role_id = role.id
		INNER JOIN role_permission ON role.id = role_permission.role_id
		INNER JOIN permission ON permission.id = role_permission.permission_id
		GROUP BY user_insert.id, user_insert.created_at, user_insert.updated_at`

	args := pgx.NamedArgs{
		"name":          u.Name,
		"email":         u.Email,
		"password_hash": u.Password.hash,
		"role":          userInactiveRole,
	}

	err := r.DB.QueryRow(ctx, query, args).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.Scopes)
	if err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "user_email_key" (SQLSTATE 23505)`:
			return errDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (r *userRepository) FindById(ctx context.Context, id string) (*user, error) {
	u := user{}
	query := `
		SELECT u.id, u.name, u.email, u.activated, u.created_at, u.updated_at, JSON_AGG(p.slug)
		FROM "user" u
		INNER JOIN public.role r ON r.id = u.role_id
		INNER JOIN public.role_permission rp ON r.id = rp.role_id
		INNER JOIN public.permission p ON rp.permission_id = p.id
		WHERE u.id = @id
		GROUP BY u.updated_at, u.created_at, u.activated, u.email, u.name, u.id`

	args := pgx.NamedArgs{
		"id": id,
	}

	err := r.DB.QueryRow(ctx, query, args).
		Scan(&u.ID, &u.Name, &u.Email, &u.Activated, &u.CreatedAt, &u.UpdatedAt, &u.Scopes)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, errUserNotFound
		default:
			return nil, err
		}
	}

	return &u, nil
}

func (r *userRepository) Activate(ctx context.Context, usr *user) error {
	query := `
		UPDATE "user"
		SET activated = @activated, 
		    updated_at = NOW(),
		    role_id = (SELECT id FROM role WHERE slug = @role)
		WHERE id = @id AND updated_at = @updated_at AND activated = FALSE
		RETURNING updated_at`

	args := pgx.NamedArgs{
		"id":         usr.ID,
		"activated":  usr.Activated,
		"updated_at": usr.UpdatedAt,
		"role":       userActiveRole,
	}

	err := r.DB.QueryRow(ctx, query, args).Scan(&usr.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil // treat as success
		default:
			return err
		}
	}

	return nil
}
