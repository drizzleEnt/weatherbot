FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/bot/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

COPY .env .

CMD ["./server"]