package user

import (
	"context"

	"github.com/saurav11sarkar/go/internal/utils"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Profile(ctx context.Context, input ProfileInput) (*User, error) {
	user, err := s.repo.Profile(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) UpdateProfile(
	ctx context.Context,
	input UpdateProfileInput,
) (*User, error) {
	user, err := s.repo.UpdateProfile(ctx, input)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetAllUsers(
	ctx context.Context,
	query utils.Query,
) (UsersResult, error) {
	if err := validateUserListQuery(query); err != nil {
		return UsersResult{}, err
	}

	users, total, err := s.repo.GetAllUsers(ctx, query)
	if err != nil {
		return UsersResult{}, err
	}

	totalPages := total / int64(query.Limit)
	if total%int64(query.Limit) != 0 {
		totalPages++
	}

	return UsersResult{
		Users:      users,
		Total:      total,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: totalPages,
	}, nil
}

func validateUserListQuery(query utils.Query) error {
	if query.Page < 1 || query.Limit < 1 || query.Limit > 100 {
		return ErrInvalidListQuery
	}
	if query.Page-1 > int(^uint(0)>>1)/query.Limit {
		return ErrInvalidListQuery
	}
	switch query.SortBy {
	case "", "name", "email", "role", "status", "createdAt", "updatedAt", "created_at", "updated_at":
	default:
		return ErrInvalidListQuery
	}
	if query.SortOrder != "asc" && query.SortOrder != "desc" {
		return ErrInvalidListQuery
	}
	for field, value := range query.Filters {
		switch field {
		case "name", "email", "photo":
		case "role":
			if value != "" && value != "user" && value != "admin" {
				return ErrInvalidListQuery
			}
		case "status":
			if value != "" && value != "active" && value != "inactive" {
				return ErrInvalidListQuery
			}
		default:
			return ErrInvalidListQuery
		}
	}
	return nil
}
