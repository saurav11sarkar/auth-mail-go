package categories

import (
	"net/http"

	"github.com/saurav11sarkar/go/internal/categories/dto"
	"github.com/saurav11sarkar/go/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CategoryRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		writeError(w, err)
		return
	}

	category, err := h.service.Create(r.Context(), CategoryInput{
		Name: req.Name,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	utils.JSON(w, http.StatusCreated, "Category created successfully", categoryResponse(category))
}

func (h *Handler) GetAllCategories(
	w http.ResponseWriter,
	r *http.Request,
) {
	q := utils.NewQuery(r.URL.Query())

	result, err := h.service.GetAllCategories(r.Context(), q)
	if err != nil {
		writeError(w, err)
		return
	}

	categories := make([]dto.CategoryResponse, 0, len(result.Categories))

	for _, category := range result.Categories {
		categories = append(categories, categoryResponse(category))
	}

	response := dto.GetCategoriesResponse{
		Categories: categories,
		Total:      result.Total,
		Page:       result.Page,
		Limit:      result.Limit,
		TotalPages: result.TotalPages,
	}

	utils.JSON(
		w,
		http.StatusOK,
		"Categories fetched successfully",
		response,
	)
}
