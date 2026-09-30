package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/saurav11sarkar/go/internal/config"
	"github.com/saurav11sarkar/go/internal/email"
	"github.com/saurav11sarkar/go/internal/utils"
)

type Service struct {
	repo  *Reposiroty
	cfg   config.Config
	email *email.Email
}

func NewService(repo *Reposiroty, cfg config.Config) *Service {
	return &Service{
		repo:  repo,
		cfg:   cfg,
		email: email.NewEmail(cfg),
	}
}

func (s *Service) Profile(ctx context.Context, id string) (*User, error) {
	user, err := s.repo.Profile(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "User not found")
		}
		return nil, err
	}
	return user, nil
}
