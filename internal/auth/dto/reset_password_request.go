package dto

type ResetPasswordRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	OTP      string `json:"otp" validate:"required,len=6,numeric"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}
