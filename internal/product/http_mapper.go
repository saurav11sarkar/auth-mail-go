package product

import (
	"errors"
	"net/http"

	"github.com/saurav11sarkar/go/internal/product/dto"
	"github.com/saurav11sarkar/go/internal/utils"
)

func productResponse(p Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		CreatedBy:   p.CreatedBy,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		PriceMinor:  p.PriceMinor,
		Currency:    p.Currency,
		Stock:       p.Stock,
		Images:      p.Images,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrSlugAlreadyExists):
		err = utils.NewAppError(http.StatusConflict, "SLUG_ALREADY_EXISTS", "Product slug already exists")
	case errors.Is(err, ErrCategoryNotFound):
		err = utils.NewAppError(http.StatusBadRequest, "INVALID_CATEGORY", "Category does not exist")
	case errors.Is(err, ErrInvalidSlug):
		err = utils.NewAppError(http.StatusBadRequest, "INVALID_SLUG", "Invalid product name for slug generation")
	case errors.Is(err, ErrProductNotFound):
		err = utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Product not found")
	}
	utils.HandlerError(w, err)
}
