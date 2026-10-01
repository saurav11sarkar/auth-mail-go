package dto

type RefreshRequestDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=8,max=255"`
}
