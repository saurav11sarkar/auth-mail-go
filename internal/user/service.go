package user

import "context"

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
