package handler

import (
	"fmt"
	"net/http"

	"report-service/service"

	"github.com/gin-gonic/gin"
)

// ReportHandler xử lý HTTP requests cho reports
type ReportHandler struct {
	reportService *service.ReportService
}

// NewReportHandler tạo instance mới của ReportHandler
func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// GenerateReport xử lý request tạo báo cáo
func (h *ReportHandler) GenerateReport(c *gin.Context) {
	var req struct {
		Type     string `json:"type" binding:"required"`
		Date     string `json:"date,omitempty"`
		FromDate string `json:"from_date,omitempty"`
		ToDate   string `json:"to_date,omitempty"`
		Format   string `json:"format,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Generate report based on type
	var csvData string
	var err error
	var filename string

	switch req.Type {
	case "daily_sales":
		csvData, err = h.reportService.DailySalesReport(req.Date)
		filename = "daily_sales_" + req.Date + ".csv"
	case "sales_range":
		csvData, err = h.reportService.SalesRangeReport(req.FromDate, req.ToDate)
		filename = "sales_range_" + req.FromDate + "_to_" + req.ToDate + ".csv"
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Unsupported report type: " + req.Type,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate report: " + err.Error(),
		})
		return
	}

	// Return CSV file for download
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "text/csv", []byte(csvData))
}

// validateRequest kiểm tra tính hợp lệ của request
func (h *ReportHandler) validateRequest(req struct {
	Type     string `json:"type" binding:"required"`
	Date     string `json:"date,omitempty"`
	FromDate string `json:"from_date,omitempty"`
	ToDate   string `json:"to_date,omitempty"`
	Format   string `json:"format,omitempty"`
}) error {
	switch req.Type {
	case "daily_sales":
		if req.Date == "" {
			return fmt.Errorf("date is required for daily_sales report")
		}
		return h.reportService.ValidateDate(req.Date)
	case "sales_range":
		if req.FromDate == "" || req.ToDate == "" {
			return fmt.Errorf("from_date and to_date are required for sales_range report")
		}
		if err := h.reportService.ValidateDate(req.FromDate); err != nil {
			return err
		}
		return h.reportService.ValidateDate(req.ToDate)
	default:
		return fmt.Errorf("unsupported report type: %s", req.Type)
	}
}
