Write-Host "Starting Consul locally..." -ForegroundColor Green
Write-Host ""

# Kiểm tra xem Consul đã được cài đặt chưa
try {
    $consulPath = Get-Command consul -ErrorAction Stop
    Write-Host "Consul found at: $($consulPath.Source)" -ForegroundColor Green
} catch {
    Write-Host "Consul not found. Please install Consul first." -ForegroundColor Red
    Write-Host "Download from: https://www.consul.io/downloads" -ForegroundColor Yellow
    Read-Host "Press Enter to continue"
    exit 1
}

Write-Host "Starting Consul agent..." -ForegroundColor Green
Write-Host ""

# Tạo thư mục data nếu chưa có
if (!(Test-Path "./consul-data")) {
    New-Item -ItemType Directory -Path "./consul-data" | Out-Null
    Write-Host "Created consul-data directory" -ForegroundColor Yellow
}

# Khởi động Consul agent
Write-Host "Consul UI will be available at: http://localhost:8500" -ForegroundColor Cyan
Write-Host "Press Ctrl+C to stop Consul" -ForegroundColor Yellow
Write-Host ""

consul agent -server -bootstrap-expect=1 -ui -client=0.0.0.0 -bind=0.0.0.0 -data-dir=./consul-data 