# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY src/ .
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server .

# Runtime stage — scratch로 공격 표면 최소화 (~5MB)
FROM scratch
COPY --from=builder /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
