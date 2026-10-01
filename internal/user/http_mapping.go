package user

import (
	"errors"
	"net/http"

	"github.com/saurav11sarkar/go/internal/user/dto"
	"github.com/saurav11sarkar/go/internal/utils"
)

func profileResponse(account *User) dto.UserProfileResponse {
	return dto.UserProfileResponse{
		ID: account.ID, Email: account.Email, Name: account.Name,
		Role: account.Role, Status: account.Status, Photo: account.Photo,
		CreatedAt: account.CreatedAt, UpdatedAt: account.UpdatedAt,
	}
}

func writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrInvalidListQuery) {
		err = utils.NewAppError(http.StatusBadRequest, "INVALID_QUERY", "Invalid pagination, sorting, or user filter")
	}
	if errors.Is(err, ErrUserNotFound) {
		err = utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "User not found")
	}
	utils.HandlerError(w, err)
}
