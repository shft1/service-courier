LOCAL_BIN := $(CURDIR)/bin
PATH := $(PATH):$(PWD)/bin

.PHONY: bin-deps
bin-deps:
	$(info installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && \
	GOBIN=$(LOCAL_BIN) go install github.com/easyp-tech/easyp/cmd/easyp@v0.7.15

.PHONY: generate
generate:
	$(info generating code...)
	@$(LOCAL_BIN)/easyp generate

.PHONY: lint
lint:
	$(info linting proto...)
	@$(LOCAL_BIN)/easyp lint --path api

.PHONY: breaking
breaking:
	$(info backwarding compatibility...)
	@$(LOCAL_BIN)/easyp breaking --against main --path api

test:
	go test -v -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

up-app:
	docker compose -f deploy/app/docker-compose.yml up -d
down-app:
	docker compose -f deploy/app/docker-compose.yml down

up-consumer:
	docker compose -f deploy/consumer/docker-compose.yml up -d
down-consumer:
	docker compose -f deploy/consumer/docker-compose.yml down

up-observability:
	docker compose -f deploy/observability/docker-compose.yml up -d

down-observability:
	docker compose -f deploy/observability/docker-compose.yml down

up-infra:
	docker compose -f ../infrastructure/docker-compose.yml up -d
down-infra:
	docker compose -f ../infrastructure/docker-compose.yml down

up-order:
	docker compose -f ../service-order/docker-compose.yaml up -d
down-order:
	docker compose -f ../service-order/docker-compose.yaml down
