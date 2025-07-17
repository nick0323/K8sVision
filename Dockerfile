# --- 构建阶段 ---
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o k8svision ./backend/main.go

# --- 运行阶段 ---
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/k8svision ./k8svision
COPY backend/config.yaml ./config.yaml
EXPOSE 8080
ENTRYPOINT ["/app/k8svision"] 