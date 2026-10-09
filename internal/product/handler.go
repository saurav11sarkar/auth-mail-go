package product

import (
	"net/http"

	"github.com/saurav11sarkar/go/internal/middlewares"
	"github.com/saurav11sarkar/go/internal/product/dto"
	"github.com/saurav11sarkar/go/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.UserContextKey).(*utils.JwtClaims)
	if !ok || claims == nil {
		utils.HandlerError(w, utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized"))
		return
	}
	var req dto.CreateProductRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		writeError(w, err)
		return
	}


	product, err := h.service.Create(r.Context(), ProductInput{
		CategoryID:  req.CategoryID,
		CreatedBy:   claims.UserID,
		Name:        req.Name,
		Description: req.Description,
		PriceMinor:  req.PriceMinor,
		Currency:    req.Currency,
		Stock:       req.Stock,
		Images:      req.Images,
		Status:      req.Status,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	utils.JSON(w, http.StatusCreated, "Product created successfully", productResponse(product))
}
