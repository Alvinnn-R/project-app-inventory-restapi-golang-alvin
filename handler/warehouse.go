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

type WarehouseHandler struct {
	WarehouseService service.WarehouseService
	Config           utils.Configuration
}

func NewWarehouseHandler(warehouseService service.WarehouseService, config utils.Configuration) WarehouseHandler {
	return WarehouseHandler{
		WarehouseService: warehouseService,
		Config:           config,
	}
}

func (h *WarehouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.WarehouseRequest
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

	warehouse := model.Warehouse{
		Name:     req.Name,
		Location: req.Location,
	}

	err = h.WarehouseService.Create(&warehouse)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "warehouse created successfully", warehouse)
}

func (h *WarehouseHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit := h.Config.Limit

	warehouses, pagination, err := h.WarehouseService.GetAllWarehouses(page, limit)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "failed to fetch warehouses: "+err.Error(), nil)
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success get data", warehouses, *pagination)
}

func (h *WarehouseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	warehouseIDstr := chi.URLParam(r, "warehouse_id")

	warehouseID, err := strconv.Atoi(warehouseIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid warehouse id", nil)
		return
	}

	warehouse, err := h.WarehouseService.GetWarehouseByID(warehouseID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get warehouse by id", warehouse)
}

func (h *WarehouseHandler) Update(w http.ResponseWriter, r *http.Request) {
	warehouseIDstr := chi.URLParam(r, "warehouse_id")

	warehouseID, err := strconv.Atoi(warehouseIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid warehouse id", nil)
		return
	}

	var req dto.WarehouseUpdateRequest
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

	warehouse := model.Warehouse{
		Name:     req.Name,
		Location: req.Location,
	}

	err = h.WarehouseService.Update(warehouseID, &warehouse)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "warehouse updated successfully", nil)
}

func (h *WarehouseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	warehouseIDstr := chi.URLParam(r, "warehouse_id")

	warehouseID, err := strconv.Atoi(warehouseIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid warehouse id", nil)
		return
	}

	err = h.WarehouseService.Delete(warehouseID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "warehouse deleted successfully", nil)
}
