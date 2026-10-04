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

func NewHandler(service *Service, cfg config.Config) *Handler {
	return &Handler{
		service: service,
		cfg:     cfg,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequestDTO
	err := utils.DecodeJSON(r, &req)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := utils.ValidateStruct(req); err != nil {
		writeError(w, err)
		return
	}
	user, err := h.service.Create(r.Context(), RegisterInput{Name: req.Name, Email: req.Email, Password: req.Password})
	if err != nil {
		writeError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, "User created successfully", registerResponse(user))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequestDTO
	err := utils.DecodeJSON(r, &req)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := utils.ValidateStruct(req); err != nil {
		writeError(w, err)
		return
	}
	user, err := h.service.Login(r.Context(), LoginInput{Email: req.Email, Password: req.Password})
	if err != nil {
		writeError(w, err)
		return
	}

	h.setRefreshCookie(w, user.RefreshToken)

	utils.JSON(w, http.StatusOK, "User login successfully", loginResponse(user))
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var token string
	// ১. যদি ব্রাউজার থেকে রিকোয়েস্ট আসে এবং কুকি থাকে:
	if cookie, err := r.Cookie("refresh_token"); err == nil && cookie.Value != "" {
		token = cookie.Value
	} else {
		// ২. যদি কুকি না থাকে (যেমন মোবাইল অ্যাপ), তবে JSON বডি ডিকোড এবং ভ্যালিডেট করতে হবে:
		var req dto.RefreshRequestDTO
		if err := utils.DecodeJSON(r, &req); err != nil {
			writeError(w, err)
			return
		}
		if err := utils.ValidateStruct(req); err != nil {
			writeError(w, err) // খালি {} পাঠালে এটা ৪০০ ভ্যালিডেশন এরর রিটার্ন করবে
			return
		}
		token = req.RefreshToken
	}

	user, err := h.service.Refresh(r.Context(), RefreshInput{RefreshToken: token})
	if err != nil {
		writeError(w, err)
		return
	}

	h.setRefreshCookie(w, user.RefreshToken)
	utils.JSON(w, http.StatusOK, "Token refreshed successfully", loginResponse(user))
}


func (h *Handler) ForgetPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ForgetPasswordRequestDTO
	err := utils.DecodeJSON(r, &req)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := utils.ValidateStruct(req); err != nil {
		writeError(w, err)
		return
	}
	err = h.service.ForgetPassword(r.Context(), ForgetPasswordInput{Email: req.Email})
	if err != nil {
		writeError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, "Password reset instructions sent successfully", nil)
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordRequestDTO
	err := utils.DecodeJSON(r, &req)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := utils.ValidateStruct(req); err != nil {
		writeError(w, err)
		return
	}
	err = h.service.ResetPassword(r.Context(), ResetPasswordInput{Email: req.Email, OTP: req.OTP, Password: req.Password})
	if err != nil {
		writeError(w, err)
		return
	}

	utils.JSON(w, http.StatusOK, "Password reset successfully", nil)
}

// Both token endpoints use the same cookie policy. MaxAge is measured in seconds.
func (h *Handler) setRefreshCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/api/v1/auth",
		HttpOnly: true,
		Secure:   h.cfg.Env == "production",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((time.Duration(h.cfg.Auth.RefreshTokenDays) * 24 * time.Hour) / time.Second),
	})
}
