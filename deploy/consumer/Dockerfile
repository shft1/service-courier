FROM golang:1.24 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o consumer ./cmd/consumer/main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /app/consumer /service-consumer
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/service-consumer"]