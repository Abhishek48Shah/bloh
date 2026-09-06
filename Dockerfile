FROM golang:1.27 AS builder
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/api

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
CMD ["./server"]

