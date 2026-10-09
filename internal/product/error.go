package product

import "errors"

var (
	ErrSlugAlreadyExists = errors.New("Product slug already exists")
	ErrProductNotFound   = errors.New("Product not found")
	ErrInvalidStatus     = errors.New("Invalid status")
	ErrCategoryNotFound  = errors.New("Category not found")
	ErrInvalidSlug       = errors.New("Invalid slug")
)
