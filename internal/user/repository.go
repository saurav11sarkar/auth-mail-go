package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/saurav11sarkar/go/internal/utils"
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

func (r *Repository) UpdateProfile(ctx context.Context, input UpdateProfileInput) (*User, error) {
	var user User
	err := r.db.QueryRow(ctx, `
		UPDATE users
		SET name=COALESCE($1, name), role=COALESCE($2, role),
		    status=COALESCE($3, status), updated_at=NOW()
		WHERE id=$4
		RETURNING id, name, email, role, status, photo, created_at, updated_at`,
		input.Name, input.Role, input.Status, input.UserID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Status, &user.Photo, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("update profile: %w", err)
	}
	return &user, nil
}

// GetAllUsers returns the requested page and the total matching user count.
func (r *Repository) GetAllUsers(ctx context.Context, q utils.Query) ([]User, int64, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 10
	}
	if q.Page-1 > int(^uint(0)>>1)/q.Limit {
		return nil, 0, fmt.Errorf("page is too large")
	}
	where, args := userFilters(q)

	var total int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users"+where, args).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	// SQL column names cannot be parameters: choose them from a fixed allowlist.
	sortColumns := map[string]string{
		"name": "name", "email": "email", "role": "role", "status": "status",
		"createdAt": "created_at", "updatedAt": "updated_at",
		"created_at": "created_at", "updated_at": "updated_at",
	}
	sortBy := sortColumns[q.SortBy]
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := "DESC"
	if q.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query := "SELECT id, email, name, role, status, photo, created_at, updated_at FROM users" + where
	query += " ORDER BY " + sortBy + " " + sortOrder + ", id " + sortOrder
	query += " LIMIT @limit OFFSET @offset"
	args["limit"] = q.Limit
	args["offset"] = q.Offset()

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, 0, fmt.Errorf("get all users: %w", err)
	}
	defer rows.Close()

	users := make([]User, 0, q.Limit)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.Status, &user.Photo, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("read users: %w", err)
	}
	return users, total, nil
}

// Reuse the same search and filters for both COUNT and SELECT.
func userFilters(q utils.Query) (string, pgx.NamedArgs) {
	where := " WHERE 1 = 1"
	args := pgx.NamedArgs{}
	if q.Search != "" {
		search := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(q.Search)
		where += " AND (name ILIKE @search OR email ILIKE @search)"
		args["search"] = "%" + search + "%"
	}
	// Unsupported filter keys are ignored; never put client-provided keys into SQL.
	for _, column := range []string{"name", "email", "role", "status", "photo"} {
		if value := q.Filters[column]; value != "" {
			where += " AND " + column + " = @" + column
			args[column] = value
		}
	}
	return where, args
}
