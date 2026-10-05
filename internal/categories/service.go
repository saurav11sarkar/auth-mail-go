package categories

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/saurav11sarkar/go/internal/config"
	"github.com/saurav11sarkar/go/internal/utils"
)

type Service struct {
	repo *Repository
	cfg  config.Config
}

// Update regenerates the slug from the new name, just like Create.
func (s *Service) Update(ctx context.Context, id string, input CategoryInput) (Category, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return Category{}, ErrInvalidCategoryInput
	}
	name := strings.TrimSpace(input.Name)
	if len([]rune(name)) < 3 || len([]rune(name)) > 100 {
		return Category{}, ErrInvalidCategoryInput
	}
	slug := regexp.MustCompile("[^a-z0-9]+").ReplaceAllString(strings.ToLower(name), "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return Category{}, ErrInvalidCategoryInput
	}
	return s.repo.Update(ctx, Category{ID: parsedID.String(), Name: name, Slug: slug})
}

func (s *Service) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidCategoryInput
	}
	return s.repo.Delete(ctx, parsedID.String())
}

func NewService(repo *Repository, cfg config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

func (s *Service) Create(ctx context.Context, input CategoryInput) (Category, error) {
	name := strings.TrimSpace(input.Name)
	slug := strings.ToLower(name)
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	if slug == "" {
		return Category{}, ErrInvalidCategoryInput
	}
	category := Category{
		Name: name,
		Slug: slug,
	}

	result, err := s.repo.Create(ctx, category)
	if err != nil {
		return Category{}, err
	}
	return result, nil
}

func (s *Service) GetAllCategories(ctx context.Context, q utils.Query) (CategoriesResult, error) {
	if q.Page < 1 || q.Limit < 1 || q.Limit > 100 {
		return CategoriesResult{}, ErrInvalidListQuery
	}

	maxInt := int(^uint(0) >> 1)
	if q.Page-1 > maxInt/q.Limit {
		return CategoriesResult{}, ErrInvalidListQuery
	}

	q.Search = strings.TrimSpace(q.Search)
	if len([]rune(q.Search)) > 200 {
		return CategoriesResult{}, ErrInvalidListQuery
	}

	switch q.SortBy {
	case "", "name", "slug", "createdAt":
	default:
		return CategoriesResult{}, ErrInvalidListQuery
	}

	switch q.SortOrder {
	case "asc", "desc":
	default:
		return CategoriesResult{}, ErrInvalidListQuery
	}

	for key, value := range q.Filters {
		switch key {
		case "name", "slug":
			if len([]rune(value)) > 200 {
				return CategoriesResult{}, ErrInvalidListQuery
			}
		default:
			return CategoriesResult{}, ErrInvalidListQuery
		}
	}

	categories, total, err := s.repo.GetAllCategories(ctx, q)
	if err != nil {
		return CategoriesResult{}, err
	}

	totalPages := total / int64(q.Limit)
	if total%int64(q.Limit) != 0 {
		totalPages++
	}

	return CategoriesResult{
		Categories: categories,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	}, nil

}

func (s *Service) GetCategoryByID(ctx context.Context, id string) (Category, error) {
	if id == "" {
		return Category{}, ErrInvalidListQuery
	}
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return Category{}, err
	}
	return category, nil
}
