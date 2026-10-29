FROM golang:alpine AS builder

WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download

COPY . .
RUN go build -o /server cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /server ./server
COPY .env /app/.env

EXPOSE 8080

CMD ["./server"]
