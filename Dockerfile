FROM golang:1.27
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o server .
CMD ["./server"]
