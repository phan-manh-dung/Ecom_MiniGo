# Report Service

Service chuyên xử lý báo cáo và export dữ liệu ra CSV.

## Cấu trúc

```
report_service/
├── main.go              # Entry point
├── go.mod               # Dependencies
├── db/
│   └── db.go           # Database connection
├── service/
│   └── report_service.go # Business logic
├── handler/
│   └── report_handler.go # HTTP handlers
└── README.md            # Documentation
```

## Cài đặt

1. **Cài đặt dependencies:**

   ```bash
   go mod tidy
   ```

2. **Cấu hình database:**

   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=123postgres
   export DB_NAME=report_service
   ```

3. **Chạy service:**
   ```bash
   go run main.go
   ```

## API Endpoints

### 1. Health Check

```
GET /health
```

### 2. Generate Report

```
POST /api/admin/reports/generate
```

#### Request Body:

```json
{
  "type": "daily_sales",
  "date": "2025-08-14"
}
```

#### Response:

```json
{
  "success": true,
  "data": "Date,Total_Orders,Revenue,Unique_Customers\n2025-08-14,150,25000000,120",
  "filename": "daily_sales_2025-08-14.csv",
  "generated_at": "2025-08-14T15:30:00Z",
  "message": "Report generated successfully",
  "type": "daily_sales"
}
```

## Loại báo cáo hỗ trợ

### 1. Daily Sales Report

- **Type**: `daily_sales`
- **Required**: `date` (YYYY-MM-DD)
- **Data**: Tổng orders, revenue, unique customers của ngày

### 2. Sales Range Report

- **Type**: `sales_range`
- **Required**: `from_date`, `to_date` (YYYY-MM-DD)
- **Data**: Tổng orders, revenue, unique customers theo khoảng thời gian

## Cách sử dụng

### Test với curl:

```bash
# Daily sales report
curl -X POST http://localhost:70051/api/admin/reports/generate \
  -H "Content-Type: application/json" \
  -d '{
    "type": "daily_sales",
    "date": "2025-08-14"
  }'

# Sales range report
curl -X POST http://localhost:70051/api/admin/reports/generate \
  -H "Content-Type: application/json" \
  -d '{
    "type": "sales_range",
    "from_date": "2025-08-01",
    "to_date": "2025-08-14"
  }'
```

### Frontend Integration:

```javascript
// Gọi API và download CSV
fetch("/api/admin/reports/generate", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    type: "daily_sales",
    date: "2025-08-14",
  }),
})
  .then((response) => response.json())
  .then((data) => {
    if (data.success) {
      // Tạo và download CSV file
      const blob = new Blob([data.data], { type: "text/csv" });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = data.filename;
      a.click();
    }
  });
```

## TODO

- [ ] Kết nối database thật để query orders
- [ ] Thêm authentication middleware
- [ ] Thêm các loại báo cáo khác
- [ ] Cache reports để tăng performance
- [ ] Background job cho reports lớn
- [ ] Export PDF, Excel
- [ ] Email reports
