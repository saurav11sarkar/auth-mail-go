package categories

import (
	"errors"
)

var (
	ErrInvalidCategoryInput  = errors.New("invalid category input")
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrInvalidListQuery      = errors.New("invalid list query")
	ErrCategoryNotFound      = errors.New("category not found")
)
