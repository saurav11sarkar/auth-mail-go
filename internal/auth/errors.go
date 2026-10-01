package auth

import "errors"

var (
	ErrEmailExists         = errors.New("email already registered")
	ErrAccountNotFound     = errors.New("account not found")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidOTP          = errors.New("invalid or expired OTP")
)
