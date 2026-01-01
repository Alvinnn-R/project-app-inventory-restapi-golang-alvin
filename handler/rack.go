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

type RackHandler struct {
	RackService service.RackService
	Config      utils.Configuration
}

func NewRackHandler(rackService service.RackService, config utils.Configuration) RackHandler {
	return RackHandler{
		RackService: rackService,
		Config:      config,
	}
}

func (h *RackHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.RackRequest
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

	rack := model.Rack{
		WarehouseID: req.WarehouseID,
		Code:        req.Code,
		Description: description,
	}

	// Create rack service
	err = h.RackService.Create(&rack)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "rack created successfully", rack)
}

func (h *RackHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	// Check if filtering by warehouse
	warehouseIDStr := r.URL.Query().Get("warehouse_id")

	// Config limit pagination
	limit := h.Config.Limit

	if warehouseIDStr != "" {
		// Get racks by warehouse
		warehouseID, err := strconv.Atoi(warehouseIDStr)
		if err != nil {
			utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid warehouse_id", nil)
			return
		}

		racks, pagination, err := h.RackService.GetRacksByWarehouse(warehouseID, page, limit)
		if err != nil {
			utils.ResponseBadRequest(w, http.StatusInternalServerError, "failed to fetch racks: "+err.Error(), nil)
			return
		}

		utils.ResponsePagination(w, http.StatusOK, "success get data", racks, *pagination)
		return
	}

	// Get all racks
	racks, pagination, err := h.RackService.GetAllRacks(page, limit)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "failed to fetch racks: "+err.Error(), nil)
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success get data", racks, *pagination)
}

func (h *RackHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	rackIDstr := chi.URLParam(r, "rack_id")

	rackID, err := strconv.Atoi(rackIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid rack id", nil)
		return
	}

	rack, err := h.RackService.GetRackByID(rackID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get rack by id", rack)
}

func (h *RackHandler) Update(w http.ResponseWriter, r *http.Request) {
	rackIDstr := chi.URLParam(r, "rack_id")

	rackID, err := strconv.Atoi(rackIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid rack id", nil)
		return
	}

	var req dto.RackUpdateRequest
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

	rack := model.Rack{
		WarehouseID: req.WarehouseID,
		Code:        req.Code,
		Description: description,
	}

	err = h.RackService.Update(rackID, &rack)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "rack updated successfully", nil)
}

func (h *RackHandler) Delete(w http.ResponseWriter, r *http.Request) {
	rackIDstr := chi.URLParam(r, "rack_id")

	rackID, err := strconv.Atoi(rackIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid rack id", nil)
		return
	}

	err = h.RackService.Delete(rackID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "rack deleted successfully", nil)
}
