package product

type ProductInput struct {
	CategoryID  string
	CreatedBy   string
	Name        string
	Description string
	PriceMinor  int64
	Currency    string
	Stock       int
	Images      []string
	Status      string
}
