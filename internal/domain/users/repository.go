package users

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const userInactiveRole = "user-inactive"

type Repository interface {
	Create(ctx context.Context, u *user) error
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

	err := r.DB.QueryRow(ctx, query, args).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.Scopes)
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
