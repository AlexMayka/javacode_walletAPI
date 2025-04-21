# Migration stage
FROM golang:1.24 AS goose-base

WORKDIR /tmp/goose
RUN go mod init goose-empty
RUN CGO_ENABLED=0 GOOS=linux go install github.com/pressly/goose/v3/cmd/goose@latest


# Build stage: compile app and generate Swagger docs
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go -o docs

RUN CGO_ENABLED=0 GOOS=linux go build -o wallet-app ./cmd/

# Final stage: minimal runtime image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/wallet-app .
COPY --from=goose-base /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/migrations /migrations

COPY migrations .
CMD ["./wallet-app"]