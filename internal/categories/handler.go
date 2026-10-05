package categories

import (
	"net/http"

	"github.com/google/uuid"
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

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.CategoryRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}
	if err := utils.ValidateStruct(req); err != nil {
		writeError(w, err)
		return
	}
	category, err := h.service.Update(r.Context(), r.PathValue("id"), CategoryInput{Name: req.Name})
	if err != nil {
		writeError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, "Category updated successfully", categoryResponse(category))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetCategoryByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		writeError(w, ErrInvalidCategoryInput)
		return
	}

	category, err := h.service.GetCategoryByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	utils.JSON(
		w,
		http.StatusOK,
		"Category fetched successfully",
		categoryResponse(category),
	)
}
