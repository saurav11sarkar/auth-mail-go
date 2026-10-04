package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *Repository) Create(ctx context.Context, user Account) (Account, error) {
	err := r.db.QueryRow(ctx, `INSERT INTO users(name,email,password) VALUES($1,$2,$3)
		RETURNING id, role, status, created_at, updated_at`, user.Name, user.Email, user.Password).
		Scan(&user.ID, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return Account{}, accountRepositoryError(err)
	}
	return user, err
}

func (r *Repository) GetEmail(ctx context.Context, email string) (Account, error) {
	var user Account
	err := r.db.QueryRow(ctx, `SELECT id, name, email, password, role, status, created_at, updated_at FROM users WHERE email=$1`, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return Account{}, accountRepositoryError(err)
	}
	return user, err
}

func (r *Repository) GetUserID(ctx context.Context, userID string) (Account, error) {
	var user Account
	err := r.db.QueryRow(ctx, `SELECT id, name, email, password, role, status, created_at, updated_at FROM users WHERE id=$1`, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return Account{}, accountRepositoryError(err)
	}
	return user, err
}

func (r *Repository) UpdateOTP(ctx context.Context, email, otp string, expireAt time.Time) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET otp=$1, otp_expires_at=$2 WHERE email=$3`, otp, expireAt, email)
	if err != nil {
		return fmt.Errorf("update password reset code: %w", err)
	}
	return nil
}

func (r *Repository) ResetPassword(ctx context.Context, email, otp, hashedPassword string) (bool, error) {
	tag, err := r.db.Exec(ctx, `UPDATE users SET password=$1, otp=NULL, otp_expires_at=NULL, updated_at=NOW()
		WHERE email=$2 AND otp=$3 AND otp_expires_at > NOW()`, hashedPassword, email, otp)
	if err != nil {
		return false, fmt.Errorf("reset password: %w", err)
	}
	return tag.RowsAffected() > 0, nil
}

// Translate database-specific lookup errors at the repository boundary.
func accountRepositoryError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
		return ErrEmailExists
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrAccountNotFound
	}
	return fmt.Errorf("account repository: %w", err)
}
