package dto

type ProfileRequestDTO struct {
	Name   string `json:"name" validate:"required,min=2,max=255"`
	Role   string `json:"role" validate:"required,oneof=user admin"`
	Status string `json:"status" validate:"required,oneof=active inactive"`
}
