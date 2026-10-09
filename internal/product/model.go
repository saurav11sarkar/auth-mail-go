package product

import "time"

type Product struct {
	ID          string
	CategoryID  string
	CreatedBy   string
	Name        string
	Slug        string
	Description string
	PriceMinor  int64 
	Currency    string
	Status      string
	Images      []string
	Stock       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
