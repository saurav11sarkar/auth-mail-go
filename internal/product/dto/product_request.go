package dto

type CreateProductRequest struct {
	CategoryID  string   `json:"categoryId" validate:"required,uuid"`
	Name        string   `json:"name" validate:"required,min=3,max=200"`
	Description string   `json:"description" validate:"max=3000"`
	PriceMinor  int64    `json:"priceMinor" validate:"required,gte=0"`
	Currency    string   `json:"currency" validate:"required,len=3"`
	Stock       int      `json:"stock" validate:"gte=0"`
	Images      []string `json:"images"`
	Status      string   `json:"status" validate:"omitempty,oneof=draft active archived"`
}
