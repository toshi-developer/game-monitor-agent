# ビルドステージ
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o agent main.go

# 実行ステージ
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/agent .
# 実行
CMD ["./agent"]