# Ecom_MiniGo

## Cấu trúc thư mục

- `api-gateway/`: Entry point, chuyển đổi HTTP ↔ gRPC, routing, auth...
- `user_service/`, `product_service/`, `order_service/`: Các microservice chính.
- `proto/`: Chứa các file định nghĩa proto cho từng domain.
- `proto/generated/`: Chứa code generate từ các file proto (dùng chung cho các service).
- `config/`: Chứa file cấu hình (yaml/json...)
- `scripts/`: Chứa các script build/generate proto/deploy.
- `.env.example`: Mẫu biến môi trường.

## Hướng dẫn build & run

```bash
# Cài đặt các dependencies
 go mod tidy

# Generate code từ proto (ví dụ)
 scripts/gen_proto.sh

# Chạy từng service
 go run ./user_service/main.go
 go run ./product_service/main.go
 go run ./order_service/main.go
 go run ./api-gateway/main.go
```

## Cấu hình

- Copy `.env.example` thành `.env` và chỉnh sửa các biến phù hợp.
- Các file cấu hình chi tiết nằm trong thư mục `config/`.
