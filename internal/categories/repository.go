package categories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/saurav11sarkar/go/internal/utils"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, category Category) (Category, error) {
	err := r.db.QueryRow(ctx, `insert into categories(name,slug) values($1,$2) returning id,name,slug,created_at,updated_at`, category.Name, category.Slug).Scan(&category.ID, &category.Name, &category.Slug, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "categories_slug_key" {
			return Category{}, ErrCategoryAlreadyExists
		}
		return Category{}, err
	}
	return category, nil
}

func (r *Repository) GetAllCategories(ctx context.Context, q utils.Query) ([]Category, int64, error) {
	where := " WHERE TRUE"
	args := pgx.NamedArgs{}

	if q.Search != "" {
		search := strings.NewReplacer(
			`\`, `\\`,
			`%`, `\%`,
			`_`, `\_`,
		).Replace(q.Search)
		where += " AND (name ILIKE @search OR slug ILIKE @search)"
		args["search"] = "%" + search + "%"
	}
	for _, colum := range []string{"name", "slug"} {
		if value := q.Filters[colum]; value != "" {
			where += " and " + colum + " = @" + colum
			args[colum] = value
		}
	}

	var total int64
	err := r.db.QueryRow(ctx, "select count(*) from categories"+where, args).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count category %w", err)
	}

	sortClomns := map[string]string{
		"name":      "name",
		"slug":      "slug",
		"createdAt": "created_at",
	}

	sortBy := sortClomns[q.SortBy]
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := q.SortOrder
	if q.SortOrder == "" {
		sortOrder = "desc"
	}
	query := `select id,name,slug,created_at,updated_at from categories` + where + " order by " + sortBy + " " + sortOrder + ", id asc" + " LIMIT @limit OFFSET @offset"
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, 0, fmt.Errorf("get categories: %w", err)
	}
	defer rows.Close()
	categories := make([]Category, 0)

	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan category: %w", err)
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate categories: %w", err)
	}
	return categories, total, nil
}
