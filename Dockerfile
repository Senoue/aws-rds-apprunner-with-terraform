FROM golang:1.23.2 AS builder

WORKDIR /app

COPY . /app/

RUN go mod download && go build -o main ./cmd/main.go

FROM alpine:3.21.1

WORKDIR /app

COPY --from=builder /app/main /app/main

RUN chmod +x /app/main

CMD ["/app/main"]
