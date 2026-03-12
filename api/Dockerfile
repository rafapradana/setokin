# Development stage
FROM golang:1.25-alpine AS development
WORKDIR /app
RUN apk add --no-cache git curl
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["air", "-c", ".air.toml"]

# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# -s -w removes symbol table and DWARF debugging information to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/api

# Production stage
FROM alpine:latest AS production
RUN apk --no-cache add ca-certificates curl
# Create a non-privileged user
RUN addgroup -S setokin && adduser -S setokin -G setokin
WORKDIR /home/setokin
COPY --from=builder /app/main .
RUN chown setokin:setokin main

USER setokin

EXPOSE 8080
CMD ["./main"]
