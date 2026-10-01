package dto

type ForgetPasswordRequestDTO struct {
	Email string `json:"email" validate:"required,email"`
}
