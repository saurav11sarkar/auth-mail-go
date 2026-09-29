package dto

type ForgetPasswordDTO struct {
	Email string `json:"email" validate:"required,email"`
}
