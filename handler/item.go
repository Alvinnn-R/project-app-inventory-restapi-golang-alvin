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

type ItemHandler struct {
	ItemService service.ItemService
	Config      utils.Configuration
}

func NewItemHandler(itemService service.ItemService, config utils.Configuration) ItemHandler {
	return ItemHandler{
		ItemService: itemService,
		Config:      config,
	}
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ItemRequest
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
	item := model.Item{
		SKU:          req.SKU,
		Name:         req.Name,
		CategoryID:   req.CategoryID,
		RackID:       req.RackID,
		Stock:        req.Stock,
		MinimumStock: req.MinimumStock,
		Price:        req.Price,
	}

	// Create item service
	err = h.ItemService.Create(&item)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "item created successfully", item)
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	// Config limit pagination
	limit := h.Config.Limit

	// Get data items from service
	items, pagination, err := h.ItemService.GetAllItems(page, limit)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "failed to fetch items: "+err.Error(), nil)
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success get data", items, *pagination)
}

func (h *ItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	itemIDstr := chi.URLParam(r, "item_id")

	itemID, err := strconv.Atoi(itemIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid item id", nil)
		return
	}

	item, err := h.ItemService.GetItemByID(itemID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get item by id", item)
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	itemIDstr := chi.URLParam(r, "item_id")

	itemID, err := strconv.Atoi(itemIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid item id", nil)
		return
	}

	var req dto.ItemUpdateRequest
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
	item := model.Item{
		SKU:          req.SKU,
		Name:         req.Name,
		CategoryID:   req.CategoryID,
		RackID:       req.RackID,
		Stock:        req.Stock,
		MinimumStock: req.MinimumStock,
		Price:        req.Price,
	}

	err = h.ItemService.Update(itemID, &item)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "item updated successfully", nil)
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	itemIDstr := chi.URLParam(r, "item_id")

	itemID, err := strconv.Atoi(itemIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid item id", nil)
		return
	}

	err = h.ItemService.Delete(itemID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "item deleted successfully", nil)
}
