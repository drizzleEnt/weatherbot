FROM golang:1.23.7 AS builder
WORKDIR /app
RUN apk add --no-cache git
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/bot/main.go

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /app/.env .
COPY --from=builder /app/server .
EXPOSE 8443
CMD ["./server"]