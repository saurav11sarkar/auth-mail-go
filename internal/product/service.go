package product

import (
	"context"
	"regexp"
	"strings"
)

// রেগেক্স প্যাকেজ লেভেলে একবার কম্পাইল করা হলো (High Performance)
var slugRegex = regexp.MustCompile(`[^a-z0-9]+`)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input ProductInput) (Product, error) {
	name := strings.TrimSpace(input.Name)
	description := strings.TrimSpace(input.Description)

	slug := strings.ToLower(name)
	slug = slugRegex.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	if slug == "" {
		return Product{}, ErrInvalidSlug
	}

	status := input.Status
	if status == "" {
		status = "draft"
	}

	images := input.Images
	if images == nil {
		images = []string{}
	}

	currency := strings.ToUpper(strings.TrimSpace(input.Currency))
	if currency == "" {
		currency = "BDT"
	}

	data, err := s.repo.Create(ctx, Product{
		CategoryID:  input.CategoryID,
		CreatedBy:   input.CreatedBy,
		Name:        name,
		Slug:        slug,
		Description: description,
		PriceMinor:  input.PriceMinor,
		Currency:    currency,
		Status:      status,
		Images:      images,
		Stock:       input.Stock,
	})
	if err != nil {
		return Product{}, err
	}

	return data, nil
}
