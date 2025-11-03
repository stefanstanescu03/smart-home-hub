FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY ./ ./

RUN go build -o client .

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/client .

CMD ["./client"]