@echo off
echo Starting Consul locally...
echo.

REM Kiểm tra xem Consul đã được cài đặt chưa
where consul >nul 2>nul
if %errorlevel% neq 0 (
    echo Consul not found. Please install Consul first.
    echo Download from: https://www.consul.io/downloads
    pause
    exit /b 1
)

echo Consul found. Starting Consul agent...
echo.

REM Khởi động Consul agent
consul agent -server -bootstrap-expect=1 -ui -client=0.0.0.0 -bind=0.0.0.0 -data-dir=./consul-data

echo.
echo Consul started successfully!
echo UI available at: http://localhost:8500
echo.
pause 