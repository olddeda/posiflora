FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /server .
COPY migrations/ migrations/
COPY locales/ locales/
EXPOSE 8080
CMD ["./server"]
