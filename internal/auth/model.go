package auth

import "time"

type Account struct {
	ID           string
	Email        string
	Password     string `json:"-"`
	Name         string
	Role         string
	Status       string
	Photo        string
	OTP          string `json:"-"`
	OTPExpiresAt time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
