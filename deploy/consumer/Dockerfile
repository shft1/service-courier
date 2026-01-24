FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o consumer ./cmd/consumer/main.go

FROM curlimages/curl:8.17.0
WORKDIR /
COPY --from=builder /app/consumer /service-consumer
EXPOSE 8080
USER nobody:nobody
ENTRYPOINT ["/service-consumer"]