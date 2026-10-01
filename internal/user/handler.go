package user

import (
	"net/http"

	"github.com/saurav11sarkar/go/internal/middlewares"
	"github.com/saurav11sarkar/go/internal/user/dto"
	"github.com/saurav11sarkar/go/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	userData, ok := r.Context().Value(middlewares.UserContextKey).(*utils.JwtClaims)
	if !ok || userData == nil {
		utils.HandlerError(w, utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized"))
		return
	}
	user, err := h.service.Profile(r.Context(), ProfileInput{UserID: userData.UserID})
	if err != nil {
		writeError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, "Profile fetched successfully", profileResponse(user))
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.UserContextKey).(*utils.JwtClaims)
	if !ok || claims == nil {
		utils.HandlerError(w, utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized"))
		return
	}
	var req dto.ProfileRequestDTO
	err := utils.DecodeJSON(r, &req)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := utils.ValidateStruct(req); err != nil {
		writeError(w, err)
		return
	}
	user, err := h.service.UpdateProfile(r.Context(), UpdateProfileInput{
		UserID: claims.UserID,
		Name:   req.Name,
		Role:   req.Role,
		Status: req.Status,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, "Profile updated successfully", profileResponse(user))
}
