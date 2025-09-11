FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o stub-service main.go

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/stub-service .

EXPOSE 9000

CMD ["./stub-service"]