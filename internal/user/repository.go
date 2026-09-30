package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Reposiroty struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Reposiroty {
	return &Reposiroty{
		db: db,
	}
}

func (r *Reposiroty) Profile(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.db.QueryRow(ctx, `
        SELECT id, email, name, role, status, photo, created_at, updated_at
        FROM users WHERE id = $1
    `, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.Status,
		&user.Photo,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
