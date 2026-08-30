FROM golang:1.27
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/api/
CMD ["./server"]
