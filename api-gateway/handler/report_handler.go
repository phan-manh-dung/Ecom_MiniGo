package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ReportHandler xử lý HTTP requests cho reports
type ReportHandler struct {
	reportServiceURL string
}

// NewReportHandler tạo instance mới của ReportHandler
func NewReportHandler() *ReportHandler {
	reportServiceURL := "http://localhost:30051"
	return &ReportHandler{
		reportServiceURL: reportServiceURL,
	}
}

// GenerateReport forward request đến Report Service
func (h *ReportHandler) GenerateReport(c *gin.Context) {
	// Đọc request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to read request body: " + err.Error(),
		})
		return
	}

	// Forward request đến Report Service
	url := h.reportServiceURL + "/api/admin/reports/generate"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to connect to Report Service: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// Đọc response từ Report Service
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to read response from Report Service: " + err.Error(),
		})
		return
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to parse response from Report Service: " + err.Error(),
		})
		return
	}

	// Forward response status và body
	c.JSON(resp.StatusCode, response)
}
