package dto

import "time"

// RegisterResponse preserves the existing registration JSON contract.
type RegisterResponse struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	Photo        string    `json:"photo,omitempty"`
	OTP          string    `json:"otp,omitempty"`
	OTPExpiresAt time.Time `json:"otp_expires_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
