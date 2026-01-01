package handler

import (
	"encoding/json"
	"net/http"
	"project-app-inventory/dto"
	"project-app-inventory/model"
	"project-app-inventory/service"
	"project-app-inventory/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CategoryHandler struct {
	CategoryService service.CategoryService
	Config          utils.Configuration
}

func NewCategoryHandler(categoryService service.CategoryService, config utils.Configuration) CategoryHandler {
	return CategoryHandler{
		CategoryService: categoryService,
		Config:          config,
	}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid request data", nil)
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Parse to model
	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	category := model.Category{
		Name:        req.Name,
		Description: description,
	}

	// Create category service
	err = h.CategoryService.Create(&category)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "category created successfully", category)
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	// Config limit pagination
	limit := h.Config.Limit

	// Get data categories from service
	categories, pagination, err := h.CategoryService.GetAllCategories(page, limit)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "failed to fetch categories: "+err.Error(), nil)
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success get data", categories, *pagination)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	categoryIDstr := chi.URLParam(r, "category_id")

	categoryID, err := strconv.Atoi(categoryIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid category id", nil)
		return
	}

	category, err := h.CategoryService.GetCategoryByID(categoryID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get category by id", category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	categoryIDstr := chi.URLParam(r, "category_id")

	categoryID, err := strconv.Atoi(categoryIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid category id", nil)
		return
	}

	var req dto.CategoryUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid request data", nil)
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Parse to model - only set fields that are provided
	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	category := model.Category{
		Name:        req.Name,
		Description: description,
	}

	err = h.CategoryService.Update(categoryID, &category)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "category updated successfully", nil)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	categoryIDstr := chi.URLParam(r, "category_id")

	categoryID, err := strconv.Atoi(categoryIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid category id", nil)
		return
	}

	err = h.CategoryService.Delete(categoryID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "category deleted successfully", nil)
}
