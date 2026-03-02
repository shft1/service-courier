FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service ./cmd/app/main.go

FROM curlimages/curl:8.17.0
WORKDIR /
COPY --from=builder /app/service /service-courier
EXPOSE 8080
USER nobody:nobody
ENTRYPOINT ["/github.com/shft1/service-courier"]