# 🛒 Ecom_MiniGo - Microservices E-commerce Platform

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![gRPC](https://img.shields.io/badge/gRPC-Protocol%20Buffers-green.svg)](https://grpc.io/)
[![Docker](https://img.shields.io/badge/Docker-Compose-orange.svg)](https://docker.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-blue.svg)](https://postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-Cache-red.svg)](https://redis.io/)

## 📋 Tổng quan dự án

**Ecom_MiniGo** là một nền tảng thương mại điện tử được xây dựng theo kiến trúc **Microservices** sử dụng **Go (Golang)** và **gRPC**. Dự án áp dụng các công nghệ và pattern hiện đại để tạo ra một hệ thống scalable, maintainable và high-performance.

### 🎯 Mục tiêu dự án

- ✅ Xây dựng hệ thống e-commerce hoàn chỉnh với microservices
- ✅ Áp dụng gRPC cho inter-service communication
- ✅ Implement các best practices trong Go development
- ✅ Sử dụng Docker và Docker Compose cho deployment
- ✅ Tích hợp service discovery với Consul
- ✅ Implement caching với Redis
- ✅ Unit testing và documentation

## 🏗️ Kiến trúc hệ thống

```
┌─────────────────┐    HTTP/REST    ┌─────────────────┐
│   Frontend      │ ──────────────► │  API Gateway    │
│   (React/Vue)   │                 │   (Port 8080)   │
└─────────────────┘                 └─────────────────┘
                                              │
                                              │ gRPC
                                              ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  User Service   │    │ Product Service │    │  Order Service  │
│   (Port 50051)  │    │   (Port 60051)  │    │   (Port 40051)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │   PostgreSQL    │    │   PostgreSQL    │
│   (User DB)     │    │  (Product DB)   │    │   (Order DB)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Công nghệ sử dụng

### Backend Stack

- **Language**: Go (Golang) 1.21+
- **Framework**: Gin (HTTP framework)
- **Communication**: gRPC + Protocol Buffers
- **Database**: PostgreSQL
- **Cache**: Redis
- **Service Discovery**: Consul
- **Containerization**: Docker & Docker Compose

### Development Tools

- **Testing**: Go testing package
- **API Documentation**: OpenAPI/Swagger (planned)
- **Code Generation**: protoc-gen-go, protoc-gen-go-grpc
- **Middleware**: CORS, Authentication, Logging, Request ID

## 🛠️ Cài đặt và chạy

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL
- Redis
- Consul (optional, for service discovery)

### Quick Start

1. **Clone repository**

```bash
git clone <repository-url>
cd Ecom_MiniGo
```

2. **Cài đặt dependencies**

```bash
go mod tidy
```

3. **Generate Protocol Buffers code**

```bash
# Linux/Mac
./scripts/gen_proto.sh

# Windows
scripts\gen_proto.bat
```

4. **Chạy với Docker Compose**

```bash
docker-compose up -d
```

5. **Hoặc chạy từng service riêng lẻ**

```bash
# Terminal 1 - User Service
go run ./user_service/main.go

# Terminal 2 - Product Service
go run ./product_service/main.go

# Terminal 3 - Order Service
go run ./order_service/main.go

# Terminal 4 - API Gateway
go run ./api-gateway/main.go
```

## 🔧 Cấu hình

### Environment Variables

Copy file cấu hình mẫu và chỉnh sửa:

```bash
cp config/config.example.yaml config/config.yaml
```

### Database Configuration

```yaml
database:
  host: localhost
  port: 5432
  name: ecom_minigo
  user: postgres
  password: your_password
```

### Redis Configuration

```yaml
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
```

## 📡 API Endpoints

### API Gateway (Port 8080)

#### User Management

```
GET    /api/users              # Lấy danh sách users
POST   /api/users              # Tạo user mới
GET    /api/users/:id          # Lấy thông tin user
PUT    /api/users/:id          # Cập nhật user
DELETE /api/users/:id          # Xóa user
```

#### Product Management

```
GET    /api/products           # Lấy danh sách products
POST   /api/products           # Tạo product mới
GET    /api/products/:id       # Lấy thông tin product
PUT    /api/products/:id       # Cập nhật product
DELETE /api/products/:id       # Xóa product
```

#### Order Management

```
GET    /api/orders             # Lấy danh sách orders
POST   /api/orders             # Tạo order mới
GET    /api/orders/:id         # Lấy thông tin order
PUT    /api/orders/:id/status  # Cập nhật trạng thái order
DELETE /api/orders/:id         # Hủy order
```

## 🧪 Testing

### Chạy Unit Tests

```bash
# Chạy tất cả tests
./run_tests.sh

# Chạy test cho service cụ thể
cd user_service/service
go test -v

# Chạy test với coverage
go test -cover
```

### Test Coverage

```bash
go test -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## 🐳 Docker

### Build Images

```bash
docker-compose build
```

### Run Services

```bash
docker-compose up -d
```

### View Logs

```bash
docker-compose logs -f [service-name]
```

## 📊 Monitoring & Health Checks

### Health Check Endpoints

- API Gateway: `http://localhost:8080/health`
- User Service: `grpc://localhost:50051/health`
- Product Service: `grpc://localhost:60051/health`
- Order Service: `grpc://localhost:40051/health`

### Service Discovery

Dự án sử dụng Consul cho service discovery:

- Consul UI: `http://localhost:8500`
- Service registration tự động khi khởi động

## 🔐 Authentication & Security

### JWT Authentication

- Middleware authentication cho API Gateway
- Token-based authentication
- Role-based access control

### CORS Configuration

- Cross-origin resource sharing được cấu hình
- Support cho frontend applications

## 📈 Performance & Scalability

### Caching Strategy

- Redis caching cho user sessions
- Product inventory caching
- Order status caching

### Database Optimization

- Connection pooling
- Indexed queries
- Pagination support

## 🚀 Deployment

### Production Deployment

```bash
# Build production images
docker-compose -f docker-compose.prod.yml build

# Deploy to production
docker-compose -f docker-compose.prod.yml up -d
```

### Environment Variables

```bash
# Production environment
export ENV=production
export DB_HOST=your-db-host
export REDIS_HOST=your-redis-host
export CONSUL_ADDR=your-consul-addr
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Team

- **Developer**: [Your Name]
- **Project**: Ecom_MiniGo Microservices Platform
- **Technology Stack**: Go, gRPC, PostgreSQL, Redis, Docker

## 📞 Support

Nếu bạn gặp vấn đề hoặc có câu hỏi:

- Tạo issue trên GitHub
- Liên hệ: [phanmanhdung2k3@gmail.com]
