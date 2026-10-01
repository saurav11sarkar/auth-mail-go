package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Profile(ctx context.Context, id string) (*User, error) {
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return &user, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, userId string, user *User) (*User, error) {
	_, err := r.db.Exec(ctx, `UPDATE users SET name=$1, role=$2, status=$3, updated_at=NOW() WHERE id=$4`, user.Name, user.Role, user.Status, userId)
	if err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}
	err = r.db.QueryRow(ctx, `SELECT name,email,role,status,photo,created_at,updated_at FROM users WHERE id=$1`, userId).Scan(&user.Name, &user.Email, &user.Role, &user.Status, &user.Photo, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get updated profile: %w", err)
	}
	return user, nil
}
