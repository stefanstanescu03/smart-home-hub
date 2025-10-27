FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

RUN go build -o server .

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/server .

EXPOSE 5000 5001

CMD ["./server"]