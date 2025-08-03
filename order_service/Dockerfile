# ---------- Build stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

# Cài tzdata để hỗ trợ TimeZone
RUN apk add --no-cache tzdata

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order_service ./order_service

# ---------- Final stage ----------
FROM alpine:latest

# Cài client Postgres và chứng chỉ TLS
RUN apk --no-cache add ca-certificates postgresql-client bash tzdata

WORKDIR /root/

# Copy binary đã build
COPY --from=builder /app/order_service .

# Copy wait-for-postgres.sh
COPY wait-for-postgres.sh /wait-for-postgres.sh
RUN chmod +x /wait-for-postgres.sh

# Expose gRPC port
EXPOSE 40051

# Đảm bảo chờ Postgres sẵn sàng trước khi chạy
ENTRYPOINT ["/wait-for-postgres.sh"]
CMD ["./order_service"]
