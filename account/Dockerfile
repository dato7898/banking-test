FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o account-service cmd/account/main.go

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/account-service .

EXPOSE 50051

CMD ["./account-service"]