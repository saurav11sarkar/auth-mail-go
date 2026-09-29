package auth

import (
	"net/http"
	"time"

	"github.com/saurav11sarkar/go/internal/auth/dto"
	"github.com/saurav11sarkar/go/internal/config"
	"github.com/saurav11sarkar/go/internal/utils"
)

type Handler struct {
	service *Service
	cfg     config.Config
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
		cfg:     config.Config{},
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var registerDto dto.RegisterDTO
	err := utils.DecodeJSON(r, &registerDto)
	if err != nil {
		utils.HandlerError(w, err)
		return
	}
	if err := utils.ValidateStruct(registerDto); err != nil {
		utils.HandlerError(w, err)
		return
	}
	user, err := h.service.Create(r.Context(), registerDto)
	if err != nil {
		utils.HandlerError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, "User created successfully", user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginDto dto.LoginDTO
	err := utils.DecodeJSON(r, &loginDto)
	if err != nil {
		utils.HandlerError(w, err)
		return
	}
	if err := utils.ValidateStruct(loginDto); err != nil {
		utils.HandlerError(w, err)
		return
	}
	user, err := h.service.Login(r.Context(), loginDto)
	if err != nil {
		utils.HandlerError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    user.RefreshToken,
		Path:     "/auth",
		HttpOnly: true,
		Secure:   h.cfg.Env == "production",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Duration(h.cfg.Auth.RefreshTokenDays) * 24 * time.Hour),
	})

	utils.JSON(w, http.StatusOK, "User login successfully", user)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var refreshDto dto.RefreshDTO
	err := utils.DecodeJSON(r, &refreshDto)
	if err != nil {
		utils.HandlerError(w, err)
		return
	}
	if err := utils.ValidateStruct(refreshDto); err != nil {
		utils.HandlerError(w, err)
		return
	}
	user, err := h.service.Refresh(r.Context(), refreshDto.RefreshToken)
	if err != nil {
		utils.HandlerError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    user.RefreshToken,
		Path:     "/auth",
		HttpOnly: true,
		Secure:   h.cfg.Env == "production",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Duration(h.cfg.Auth.RefreshTokenDays) * 24 * time.Hour),
	})

	utils.JSON(w, http.StatusOK, "User login successfully", user)
}

func (h *Handler) ForgetPassword(w http.ResponseWriter, r *http.Request) {
	var forgetPasswordDto dto.ForgetPasswordDTO
	err := utils.DecodeJSON(r, &forgetPasswordDto)
	if err != nil {
		utils.HandlerError(w, err)
		return
	}
	if err := utils.ValidateStruct(forgetPasswordDto); err != nil {
		utils.HandlerError(w, err)
		return
	}
	err = h.service.ForgetPassword(r.Context(), forgetPasswordDto)
	if err != nil {
		utils.HandlerError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, "Password reset instructions sent successfully", nil)
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var resetPasswordDto dto.ResetPasswordDTO
	err := utils.DecodeJSON(r, &resetPasswordDto)
	if err != nil {
		utils.HandlerError(w, err)
		return
	}
	if err := utils.ValidateStruct(resetPasswordDto); err != nil {
		utils.HandlerError(w, err)
		return
	}
	err = h.service.ResetPassword(r.Context(), resetPasswordDto)
	if err != nil {
		utils.HandlerError(w, err)
		return
	}

	utils.JSON(w, http.StatusOK, "Password reset successfully", nil)
}
