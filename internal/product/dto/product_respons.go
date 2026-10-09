package dto

import "time"

type ProductResponse struct {
	ID          string    `json:"id"`
	CategoryID  string    `json:"categoryId"`
	CreatedBy   string    `json:"createdBy"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	PriceMinor  int64     `json:"priceMinor"`
	Currency    string    `json:"currency"`
	Stock       int       `json:"stock"`
	Images      []string  `json:"images"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
