package auth

import (
	"errors"
	"net/http"

	"github.com/saurav11sarkar/go/internal/auth/dto"
	"github.com/saurav11sarkar/go/internal/utils"
)

// Explicit mapping prevents internal credentials from reaching JSON responses.
func registerResponse(account Account) dto.RegisterResponse {
	return dto.RegisterResponse{
		ID: account.ID, Email: account.Email, Name: account.Name,
		Role: account.Role, Status: account.Status, Photo: account.Photo,
		CreatedAt: account.CreatedAt, UpdatedAt: account.UpdatedAt,
	}
}

func loginResponse(result LoginResult) dto.LoginResponse {
	account := result.Account
	return dto.LoginResponse{
		ID: account.ID, Email: account.Email, Name: account.Name,
		Role: account.Role, Status: account.Status, Photo: account.Photo,
		CreatedAt: account.CreatedAt, UpdatedAt: account.UpdatedAt,
		AccessToken: result.AccessToken, RefreshToken: result.RefreshToken,
	}
}

// HTTP status codes belong at the HTTP boundary, not in business services.
func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrEmailExists):
		err = utils.NewAppError(http.StatusConflict, "EMAIL_ALREADY_EXISTS", "Email already registered")
	case errors.Is(err, ErrInvalidCredentials):
		err = utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Invalid email or password")
	case errors.Is(err, ErrInvalidRefreshToken):
		err = utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Invalid refresh token")
	case errors.Is(err, ErrInvalidOTP):
		err = utils.NewAppError(http.StatusBadRequest, "INVALID_OR_EXPIRED_OTP", "Invalid or expired OTP")
	}
	utils.HandlerError(w, err)
}
