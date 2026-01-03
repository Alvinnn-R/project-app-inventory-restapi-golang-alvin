package handler

import (
	"net/http"
	"project-app-inventory/service"
	"project-app-inventory/utils"
)

type ReportHandler struct {
	ReportService service.ReportService
}

func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{ReportService: reportService}
}

func (h *ReportHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	report, err := h.ReportService.GetSummary()
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "failed to fetch report: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get report summary", report)
}
