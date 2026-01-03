package dto

type ReportSummaryResponse struct {
	TotalItems      int     `json:"total_items"`
	LowStockItems   int     `json:"low_stock_items"`
	TotalSales      int     `json:"total_sales"`
	TotalRevenue    float64 `json:"total_revenue"`
	ActiveUsers     int     `json:"active_users"`
	TotalCategories int     `json:"total_categories"`
	TotalWarehouses int     `json:"total_warehouses"`
}
