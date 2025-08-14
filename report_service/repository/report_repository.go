package repository

import (
	"report-service/db"
	"time"
)

// ReportRepository xử lý query database cho reports
type ReportRepository struct{}

// NewReportRepository tạo instance mới
func NewReportRepository() *ReportRepository {
	return &ReportRepository{}
}

// GetDailySalesData lấy data báo cáo theo ngày
func (r *ReportRepository) GetDailySalesData(date string) (*DailySalesData, error) {
	var data DailySalesData
	data.Date = date

	// Parse date string
	reportDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	// Query total orders và unique users với status filter
	var orderResult struct {
		TotalOrders int `json:"total_orders"`
		UniqueUsers int `json:"unique_users"`
	}

	err = db.DB.Table("orders").
		Select("COUNT(*) as total_orders, COUNT(DISTINCT user_id) as unique_users").
		Where("DATE(created_at) = ? AND status IN (?, ?, ?)", reportDate, "pending", "completed", "cancelled").
		Scan(&orderResult).Error

	if err != nil {
		return nil, err
	}

	// Query status breakdown
	var statusResult struct {
		PendingOrders   int `json:"pending_orders"`
		CompletedOrders int `json:"completed_orders"`
		CancelledOrders int `json:"cancelled_orders"`
	}

	err = db.DB.Table("orders").
		Select(`
			SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending_orders,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed_orders,
			SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END) as cancelled_orders
		`).
		Where("DATE(created_at) = ?", reportDate).
		Scan(&statusResult).Error

	if err != nil {
		return nil, err
	}

	// Query revenue từ order_details với status filter
	var revenueResult struct {
		Revenue float64 `json:"revenue"`
	}

	err = db.DB.Table("order_details od").
		Select("SUM(od.quantity * od.unit_price) as revenue").
		Joins("JOIN orders o ON o.id = od.order_id").
		Where("DATE(o.created_at) = ? AND o.status IN (?, ?, ?)", reportDate, "pending", "completed", "cancelled").
		Scan(&revenueResult).Error

	if err != nil {
		return nil, err
	}

	data.TotalOrders = orderResult.TotalOrders
	data.Revenue = revenueResult.Revenue
	data.UniqueCustomers = orderResult.UniqueUsers
	data.PendingOrders = statusResult.PendingOrders
	data.CompletedOrders = statusResult.CompletedOrders
	data.CancelledOrders = statusResult.CancelledOrders

	// Query top product với status filter
	var topProduct struct {
		ProductName string `json:"product_name"`
		Category    string `json:"category"`
		TotalSold   int    `json:"total_sold"`
	}

	err = db.DB.Table("order_details od").
		Select("p.name as product_name, p.category, SUM(od.quantity) as total_sold").
		Joins("JOIN products p ON p.id = od.product_id").
		Joins("JOIN orders o ON o.id = od.order_id").
		Where("DATE(o.created_at) = ? AND o.status IN (?, ?, ?)", reportDate, "pending", "completed", "cancelled").
		Group("od.product_id, p.name, p.category").
		Order("total_sold DESC").
		Limit(1).
		Scan(&topProduct).Error

	if err == nil {
		data.TopProduct = topProduct.ProductName
		data.TopCategory = topProduct.Category
	}

	return &data, nil
}

// GetSalesRangeData lấy data báo cáo theo khoảng thời gian
func (r *ReportRepository) GetSalesRangeData(fromDate, toDate string) ([]DailySalesData, error) {
	var dataList []DailySalesData

	// Parse dates
	from, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return nil, err
	}

	to, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return nil, err
	}

	// Query data cho từng ngày
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		dailyData, err := r.GetDailySalesData(dateStr)
		if err != nil {
			continue // Skip ngày có lỗi
		}
		dataList = append(dataList, *dailyData)
	}

	return dataList, nil
}

// DailySalesData struct chứa data báo cáo
type DailySalesData struct {
	Date            string  `json:"date"`
	TotalOrders     int     `json:"total_orders"`
	Revenue         float64 `json:"revenue"`
	UniqueCustomers int     `json:"unique_customers"`
	TopProduct      string  `json:"top_product"`
	TopCategory     string  `json:"top_category"`
	PendingOrders   int     `json:"pending_orders"`
	CompletedOrders int     `json:"completed_orders"`
	CancelledOrders int     `json:"cancelled_orders"`
}
