# ---------- Build stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

# Cài tzdata để hỗ trợ TimeZone
RUN apk add --no-cache tzdata

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o product_service ./product_service

# ---------- Final stage ----------
FROM alpine:latest

# Cài client Postgres và chứng chỉ TLS
RUN apk --no-cache add ca-certificates postgresql-client bash tzdata

WORKDIR /root/

# Copy binary đã build
COPY --from=builder /app/product_service .

# Copy wait-for-postgres script
COPY wait-for-postgres.sh /wait-for-postgres.sh
RUN chmod +x /wait-for-postgres.sh

# Expose gRPC port
EXPOSE 60051

# Sửa chỗ này:
ENTRYPOINT ["/wait-for-postgres.sh"]
CMD ["./product_service"]
