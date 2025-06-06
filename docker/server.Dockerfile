FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/bot/main.go

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /app/.env .
COPY --from=builder /app/server .
EXPOSE 8443
CMD ["./server"]