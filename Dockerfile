FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o fitness-tracker ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/fitness-tracker .
COPY .env .
EXPOSE 8080
ENTRYPOINT ["./fitness-tracker"]
