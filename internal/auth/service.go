package auth

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/saurav11sarkar/go/internal/auth/dto"
	"github.com/saurav11sarkar/go/internal/config"
	"github.com/saurav11sarkar/go/internal/email"
	"github.com/saurav11sarkar/go/internal/utils"
)

type Service struct {
	repo  *Resposistory
	cfg   config.Config
	email *email.Email
}

func NewService(repo *Resposistory, cfg config.Config) *Service {
	return &Service{
		repo:  repo,
		cfg:   cfg,
		email: email.NewEmail(cfg),
	}
}

func (s *Service) Create(ctx context.Context, registerDto dto.RegisterDTO) (Register, error) {
	email := strings.ToLower(strings.TrimSpace(registerDto.Email))
	_, err := s.repo.GetEmail(ctx, email)
	if err == nil {
		return Register{}, utils.NewAppError(http.StatusConflict, "EMAIL_ALREADY_EXISTS", "Email already registered")
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return Register{}, utils.NewAppError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error")
	}
	hash, err := utils.HashPassword(registerDto.Password)
	if err != nil {
		return Register{}, utils.NewAppError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error")
	}
	user := Register{
		ID:       utils.NewID(),
		Email:    email,
		Password: hash,
		Name:     registerDto.Name,
		Role:     "user",
		Status:   "active",
	}
	user, err = s.repo.Create(ctx, user)
	if err != nil {
		return Register{}, utils.NewAppError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error")
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, loginDto dto.LoginDTO) (Login, error) {
	user, err := s.repo.GetEmail(ctx, strings.ToLower(strings.TrimSpace(loginDto.Email)))
	if err != nil {
		return Login{}, utils.NewAppError(
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"Invalid email or password",
		)
	}

	if !utils.ComparePassword(loginDto.Password, user.Password) {
		return Login{}, utils.NewAppError(
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"Invalid email or password",
		)
	}

	accessToken, err := utils.CreateJwtToken(
		user.ID,
		user.Role,
		"access",
		s.cfg.Auth.JwtAccessSecret,
		time.Duration(s.cfg.Auth.AccessTokenMinutes)*time.Minute,
	)
	if err != nil {
		return Login{}, utils.NewAppError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			"Internal server error",
		)
	}

	refreshToken, err := utils.CreateJwtToken(
		user.ID,
		user.Role,
		"refresh",
		s.cfg.Auth.JwtRefreshSecret,
		time.Duration(s.cfg.Auth.RefreshTokenDays)*24*time.Hour,
	)
	if err != nil {
		return Login{}, utils.NewAppError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			"Internal server error",
		)
	}

	return Login{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Role:         user.Role,
		Status:       user.Status,
		Photo:        user.Photo,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (Login, error) {
	claims, err := utils.ParseJwtToken(refreshToken, s.cfg.Auth.JwtRefreshSecret, "refresh")
	if err != nil {
		return Login{}, utils.NewAppError(
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"Invalid refresh token",
		)
	}
	user, err := s.repo.GetUserID(ctx, claims.UserID)
	if err != nil {
		return Login{}, utils.NewAppError(
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"Invalid refresh token",
		)
	}

	if user.Status != "active" {
		return Login{}, utils.NewAppError(
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"Invalid refresh token",
		)
	}

	accessToken, err := utils.CreateJwtToken(
		user.ID,
		user.Role,
		"access",
		s.cfg.Auth.JwtAccessSecret,
		time.Duration(s.cfg.Auth.AccessTokenMinutes)*time.Minute,
	)
	if err != nil {
		return Login{}, utils.NewAppError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			"Internal server error",
		)
	}

	newRefreshToken, err := utils.CreateJwtToken(
		user.ID,
		user.Role,
		"refresh",
		s.cfg.Auth.JwtRefreshSecret,
		time.Duration(s.cfg.Auth.RefreshTokenDays)*24*time.Hour,
	)
	if err != nil {
		return Login{}, utils.NewAppError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			"Internal server error",
		)
	}

	return Login{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Role:         user.Role,
		Status:       user.Status,
		Photo:        user.Photo,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (s *Service) ForgetPassword(ctx context.Context, forgetPasswordDto dto.ForgetPasswordDTO) error {
	email := strings.ToLower(strings.TrimSpace(forgetPasswordDto.Email))

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

func (s *Service) ResetPassword(ctx context.Context, resetPasswordDto dto.ResetPasswordDTO) error {
	email := strings.ToLower(strings.TrimSpace(resetPasswordDto.Email))
	hashedPassword, err := utils.HashPassword(resetPasswordDto.Password)
	if err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error")
	}

	reset, err := s.repo.ResetPassword(ctx, email, resetPasswordDto.OTP, hashedPassword)
	if err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error")
	}
	if !reset {
		return utils.NewAppError(http.StatusBadRequest, "INVALID_OR_EXPIRED_OTP", "Invalid or expired OTP")
	}
	return nil
}
