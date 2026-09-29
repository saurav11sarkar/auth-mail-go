package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Resposistory struct {
	db *pgxpool.Pool
}

func NewResposistory(db *pgxpool.Pool) *Resposistory {
	return &Resposistory{
		db: db,
	}
}

func (r *Resposistory) Create(ctx context.Context, user Register) (Register, error) {
	err := r.db.QueryRow(ctx, `INSERT INTO users(id,name,email,password) VALUES($1,$2,$3,$4)
		RETURNING id, role, status, created_at, updated_at`, user.ID, user.Name, user.Email, user.Password).
		Scan(&user.ID, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return Register{}, err
	}
	return user, err
}

func (r *Resposistory) GetEmail(ctx context.Context, email string) (Register, error) {
	var user Register
	err := r.db.QueryRow(ctx, `SELECT id, name, email, password, role, status, created_at, updated_at FROM users WHERE email=$1`, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return Register{}, err
	}
	return user, err
}

func (r *Resposistory) GetUserID(ctx context.Context, userID string) (Register, error) {
	var user Register
	err := r.db.QueryRow(ctx, `SELECT id, name, email, password, role, status, created_at, updated_at FROM users WHERE id=$1`, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return Register{}, err
	}
	return user, err
}

func (r *Resposistory) UpdateOTP(ctx context.Context, email, otp string, expireAt time.Time) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET otp=$1, otp_expires_at=$2 WHERE email=$3`, otp, expireAt, email)
	return err
}

func (r *Resposistory) ResetPassword(ctx context.Context, email, otp, hashedPassword string) (bool, error) {
	tag, err := r.db.Exec(ctx, `UPDATE users SET password=$1, otp=NULL, otp_expires_at=NULL
		WHERE email=$2 AND otp=$3 AND otp_expires_at > NOW()`, hashedPassword, email, otp)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}
