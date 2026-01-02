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

type SaleHandler struct {
	SaleService service.SaleService
	Config      utils.Configuration
}

func NewSaleHandler(saleService service.SaleService, config utils.Configuration) SaleHandler {
	return SaleHandler{
		SaleService: saleService,
		Config:      config,
	}
}

func (h *SaleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.SaleRequest
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

	// Get user from context (set by auth middleware)
	userCtx := r.Context().Value("user")
	if userCtx == nil {
		utils.ResponseBadRequest(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	user, ok := userCtx.(*model.User)
	if !ok {
		utils.ResponseBadRequest(w, http.StatusUnauthorized, "invalid user context", nil)
		return
	}

	sale, err := h.SaleService.Create(user.ID, req.Items)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "sale created successfully", sale)
}

func (h *SaleHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit := h.Config.Limit

	sales, pagination, err := h.SaleService.GetAllSales(page, limit)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "failed to fetch sales: "+err.Error(), nil)
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success get data", sales, *pagination)
}

func (h *SaleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	saleIDstr := chi.URLParam(r, "sale_id")

	saleID, err := strconv.Atoi(saleIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid sale id", nil)
		return
	}

	sale, items, err := h.SaleService.GetSaleByID(saleID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	// Build response with items
	response := dto.SaleResponse{
		ID:          sale.ID,
		UserID:      sale.UserID,
		TotalAmount: sale.TotalAmount,
		CreatedAt:   sale.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   sale.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if sale.DeletedAt != nil {
		deletedAtStr := sale.DeletedAt.Format("2006-01-02 15:04:05")
		response.DeletedAt = &deletedAtStr
	}

	var saleItems []dto.SaleItemResponse
	for _, item := range items {
		saleItems = append(saleItems, dto.SaleItemResponse{
			ID:          item.ID,
			ItemID:      item.ItemID,
			Quantity:    item.Quantity,
			PriceAtSale: item.PriceAtSale,
			Subtotal:    item.Subtotal,
		})
	}
	response.Items = saleItems

	utils.ResponseSuccess(w, http.StatusOK, "success get sale by id", response)
}

func (h *SaleHandler) Update(w http.ResponseWriter, r *http.Request) {
	saleIDstr := chi.URLParam(r, "sale_id")

	saleID, err := strconv.Atoi(saleIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid sale id", nil)
		return
	}

	var req dto.SaleRequest
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

	err = h.SaleService.Update(saleID, req.Items)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "sale updated successfully", nil)
}

func (h *SaleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	saleIDstr := chi.URLParam(r, "sale_id")

	saleID, err := strconv.Atoi(saleIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid sale id", nil)
		return
	}

	err = h.SaleService.Delete(saleID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "sale deleted successfully", nil)
}
