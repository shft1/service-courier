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

up-infra:
	docker compose -f ../infrastructure/docker-compose.yml up -d
down-infra:
	docker compose -f ../infrastructure/docker-compose.yml down

up-order:
	docker compose -f ../service-order/docker-compose.yaml up -d
down-order:
	docker compose -f ../service-order/docker-compose.yaml down
