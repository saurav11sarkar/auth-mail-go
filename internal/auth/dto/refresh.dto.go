package dto

type RefreshDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=8,max=255"`
}
