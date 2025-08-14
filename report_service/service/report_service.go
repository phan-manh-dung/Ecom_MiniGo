package service

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"report-service/repository"
)

// ReportService xử lý logic tạo báo cáo
type ReportService struct {
	reportRepo *repository.ReportRepository
}

// NewReportService tạo instance mới của ReportService
func NewReportService(reportRepo *repository.ReportRepository) *ReportService {
	return &ReportService{
		reportRepo: reportRepo,
	}
}

// DailySalesReport tạo báo cáo doanh số theo ngày
func (s *ReportService) DailySalesReport(date string) (string, error) {
	// Query database thật
	data, err := s.reportRepo.GetDailySalesData(date)
	if err != nil {
		return "", fmt.Errorf("failed to get daily sales data: %v", err)
	}

	// Tạo CSV với data thật
	reportData := [][]string{
		{"Date", "Total_Orders", "Revenue", "Unique_Customers", "Top_Product", "Top_Category", "Pending_Orders", "Completed_Orders", "Cancelled_Orders"},
		{
			data.Date,
			strconv.Itoa(data.TotalOrders),
			strconv.FormatFloat(data.Revenue, 'f', 0, 64),
			strconv.Itoa(data.UniqueCustomers),
			data.TopProduct,
			data.TopCategory,
			strconv.Itoa(data.PendingOrders),
			strconv.Itoa(data.CompletedOrders),
			strconv.Itoa(data.CancelledOrders),
		},
	}

	// Tạo CSV string
	var csvBuilder strings.Builder
	writer := csv.NewWriter(&csvBuilder)

	for _, row := range reportData {
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %v", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("failed to flush CSV: %v", err)
	}

	return csvBuilder.String(), nil
}

// SalesRangeReport tạo báo cáo doanh số theo khoảng thời gian
func (s *ReportService) SalesRangeReport(fromDate, toDate string) (string, error) {
	// Query database thật
	dataList, err := s.reportRepo.GetSalesRangeData(fromDate, toDate)
	if err != nil {
		return "", fmt.Errorf("failed to get sales range data: %v", err)
	}

	// Tạo CSV với data
	reportData := [][]string{
		{"Date", "Total_Orders", "Revenue", "Unique_Customers", "Top_Product", "Top_Category", "Pending_Orders", "Completed_Orders", "Cancelled_Orders"},
	}

	// Thêm data cho từng ngày
	for _, data := range dataList {
		row := []string{
			data.Date,
			strconv.Itoa(data.TotalOrders),
			strconv.FormatFloat(data.Revenue, 'f', 0, 64),
			strconv.Itoa(data.UniqueCustomers),
			data.TopProduct,
			data.TopCategory,
			strconv.Itoa(data.PendingOrders),
			strconv.Itoa(data.CompletedOrders),
			strconv.Itoa(data.CancelledOrders),
		}
		reportData = append(reportData, row)
	}

	// Tạo CSV string
	var csvBuilder strings.Builder
	writer := csv.NewWriter(&csvBuilder)

	for _, row := range reportData {
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %v", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("failed to flush CSV: %v", err)
	}

	return csvBuilder.String(), nil
}

// ValidateDate kiểm tra format date
func (s *ReportService) ValidateDate(date string) error {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %v", err)
	}
	return nil
}
