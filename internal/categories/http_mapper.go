package categories

import (
	"errors"
	"net/http"

	"github.com/saurav11sarkar/go/internal/categories/dto"
	"github.com/saurav11sarkar/go/internal/utils"
)

func categoryResponse(category Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidCategoryInput):
		err = utils.NewAppError(http.StatusBadRequest, "INVALID_CATEGORY", "Invalid category input")
	case errors.Is(err, ErrCategoryAlreadyExists):
		err = utils.NewAppError(http.StatusConflict, "CATEGORY_ALREADY_EXISTS", "Category already exists")
	case errors.Is(err, ErrInvalidListQuery):
		err = utils.NewAppError(http.StatusBadRequest, "INVALID_LIST_QUERY", "Invalid list query")
	case errors.Is(err, ErrCategoryNotFound):
		err = utils.NewAppError(http.StatusNotFound, "CATEGORY_NOT_FOUND", "Category not found")
	}
	utils.HandlerError(w, err)
}
