package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, p Product) (Product, error) {
	query := `
		insert into products (category_id, created_by, name, slug, description, price_minor, currency, status, images, stock)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		returning id, category_id, created_by, name, slug, description, price_minor, currency, status, images, stock, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		p.CategoryID, p.CreatedBy, p.Name, p.Slug, p.Description, p.PriceMinor, p.Currency, p.Status, p.Images, p.Stock,
	).Scan(
		&p.ID, &p.CategoryID, &p.CreatedBy, &p.Name, &p.Slug, &p.Description,
		&p.PriceMinor, &p.Currency, &p.Status, &p.Images, &p.Stock, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			if pgxErr.Code == "23505" && pgxErr.ConstraintName == "products_slug_key" {
				return Product{}, ErrSlugAlreadyExists
			}
			if pgxErr.Code == "23503" && pgxErr.ConstraintName == "products_category_id_fkey" {
				return Product{}, ErrCategoryNotFound
			}
		}
		return Product{}, fmt.Errorf("failed to create product: %w", err)
	}

	return p, nil
}
