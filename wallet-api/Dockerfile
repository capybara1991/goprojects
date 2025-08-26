FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o auth-service ./cmd/server

FROM scratch
COPY --from=builder /app/auth-service /auth-service
ENTRYPOINT ["/auth-service"]