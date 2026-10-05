package dto

type CategoryRequest struct {
	Name string `json:"name" validate:"required,max=100,min=3"`
}
