package auth

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/saurav11sarkar/go/internal/config"
	"github.com/saurav11sarkar/go/internal/email"
	"github.com/saurav11sarkar/go/internal/utils"
)

type Service struct {
	repo  *Repository
	cfg   config.Config
	email *email.Email
}

func NewService(repo *Repository, cfg config.Config) *Service {
	return &Service{
		repo:  repo,
		cfg:   cfg,
		email: email.NewEmail(cfg),
	}
}

func (s *Service) Create(ctx context.Context, input RegisterInput) (Account, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	_, err := s.repo.GetEmail(ctx, email)
	if err == nil {
		return Account{}, ErrEmailExists
	}
	if !errors.Is(err, ErrAccountNotFound) {
		return Account{}, fmt.Errorf("authentication operation: %w", err)
	}
	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return Account{}, fmt.Errorf("authentication operation: %w", err)
	}
	user := Account{
		Email:    email,
		Password: hash,
		Name:     input.Name,
		Role:     "user",
		Status:   "active",
	}
	user, err = s.repo.Create(ctx, user)
	if err != nil {
		return Account{}, fmt.Errorf("authentication operation: %w", err)
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (LoginResult, error) {
	user, err := s.repo.GetEmail(ctx, strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, fmt.Errorf("login lookup: %w", err)
	}

	if !utils.ComparePassword(input.Password, user.Password) {
		return LoginResult{}, ErrInvalidCredentials
	}

	accessToken, err := utils.CreateJwtToken(
		user.ID,
		user.Role,
		"access",
		s.cfg.Auth.JwtAccessSecret,
		time.Duration(s.cfg.Auth.AccessTokenMinutes)*time.Minute,
	)
	if err != nil {
		return LoginResult{}, fmt.Errorf("issue token: %w", err)
	}

	refreshToken, err := utils.CreateJwtToken(
		user.ID,
		user.Role,
		"refresh",
		s.cfg.Auth.JwtRefreshSecret,
		time.Duration(s.cfg.Auth.RefreshTokenDays)*24*time.Hour,
	)
	if err != nil {
		return LoginResult{}, fmt.Errorf("issue token: %w", err)
	}

	return LoginResult{
		Account:      user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, input RefreshInput) (LoginResult, error) {
	claims, err := utils.ParseJwtToken(input.RefreshToken, s.cfg.Auth.JwtRefreshSecret, "refresh")
	if err != nil {
		return LoginResult{}, ErrInvalidRefreshToken
	}
	user, err := s.repo.GetUserID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			return LoginResult{}, ErrInvalidRefreshToken
		}
		return LoginResult{}, fmt.Errorf("refresh lookup: %w", err)
	}

	if user.Status != "active" {
		return LoginResult{}, ErrInvalidRefreshToken
	}

	accessToken, err := utils.CreateJwtToken(
		user.ID,
		user.Role,
		"access",
		s.cfg.Auth.JwtAccessSecret,
		time.Duration(s.cfg.Auth.AccessTokenMinutes)*time.Minute,
	)
	if err != nil {
		return LoginResult{}, fmt.Errorf("issue token: %w", err)
	}

	newRefreshToken, err := utils.CreateJwtToken(
		user.ID,
		user.Role,
		"refresh",
		s.cfg.Auth.JwtRefreshSecret,
		time.Duration(s.cfg.Auth.RefreshTokenDays)*24*time.Hour,
	)
	if err != nil {
		return LoginResult{}, fmt.Errorf("issue token: %w", err)
	}

	return LoginResult{
		Account:      user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) ForgetPassword(ctx context.Context, input ForgetPasswordInput) error {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	_, err := s.repo.GetEmail(ctx, email)
	if err != nil {
		return err
	}

	otpNumber := rand.Intn(900000) + 100000
	otp := fmt.Sprintf("%d", otpNumber)
	expireAt := time.Now().Add(15 * time.Minute)

	err = s.repo.UpdateOTP(ctx, email, otp, expireAt)
	if err != nil {
		return err
	}

	err = s.email.Send(email, "Reset Password", "<h1>Reset Password</h1><p>Your password reset code is: <strong>"+otp+"</strong></p>")
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return fmt.Errorf("authentication operation: %w", err)
	}

	reset, err := s.repo.ResetPassword(ctx, email, input.OTP, hashedPassword)
	if err != nil {
		return fmt.Errorf("authentication operation: %w", err)
	}
	if !reset {
		return ErrInvalidOTP
	}
	return nil
}
