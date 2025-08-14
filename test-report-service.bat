@echo off
echo Testing Report Service...
echo.

echo 1. Testing Daily Sales Report...
curl -X POST http://localhost:8080/api/admin/reports/generate ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer YOUR_JWT_TOKEN" ^
  -d "{\"type\": \"daily_sales\", \"date\": \"2025-08-14\"}"

echo.
echo.
echo 2. Testing Sales Range Report...
curl -X POST http://localhost:8080/api/admin/reports/generate ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer YOUR_JWT_TOKEN" ^
  -d "{\"type\": \"sales_range\", \"from_date\": \"2025-08-01\", \"to_date\": \"2025-08-14\"}"

echo.
echo.
echo Test completed!
pause 