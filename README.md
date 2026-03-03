# service-courier

## Содержание
- [0. Используемые технологии](#0-используемые-технологии)
- [1. О проекте](#1-о-проекте)
- [2. Ключевые цели проекта](#2-ключевые-цели-проекта)
- [3. Предметная область и роль сервиса](#3-предметная-область-и-роль-сервиса)
- [4. Архитектура проекта](#4-архитектура-проекта)
- [5. Реализованный функционал по этапам](#5-реализованный-функционал-по-этапам)
  - [5.1 Базовый HTTP-сервис и graceful shutdown](#51-базовый-http-сервис-и-graceful-shutdown)
  - [5.2 PostgreSQL, миграции и CRUD курьеров](#52-postgresql-миграции-и-crud-курьеров)
  - [5.3 Слоистая архитектура и dependency injection](#53-слоистая-архитектура-и-dependency-injection)
  - [5.4 Назначение/снятие курьера с заказа](#54-назначениеснятие-курьера-с-заказа)
  - [5.5 Конкурентность и фоновые воркеры](#55-конкурентность-и-фоновые-воркеры)
  - [5.6 Unit и интеграционные тесты](#56-unit-и-интеграционные-тесты)
  - [5.7 Интеграция с внешним сервисом заказов](#57-интеграция-с-внешним-сервисом-заказов)
  - [5.8 Асинхронная интеграция через Kafka](#58-асинхронная-интеграция-через-kafka)
  - [5.9 Observability: метрики, логирование, мониторинг](#59-observability-метрики-логирование-мониторинг)
  - [5.10 Rate Limiter и Retry](#510-rate-limiter-и-retry)
  - [5.11 CI/CD и качество кода](#511-cicd-и-качество-кода)
  - [5.12 Профилирование и оптимизация](#512-профилирование-и-оптимизация)
- [6. API сервиса](#6-api-сервиса)
- [7. Хранилище данных](#7-хранилище-данных)
- [8. Конфигурация](#8-конфигурация)
- [9. Локальный запуск](#9-локальный-запуск)

---

## 0. Используемые технологии
- **Go 1.24** — основной язык разработки сервиса.
- **Chi (net/http router)** — HTTP-роутинг и middleware.
- **PostgreSQL** — основное хранилище данных.
- **pgx / pgxpool** — драйвер и пул соединений с PostgreSQL.
- **Goose** — миграции схемы базы данных.
- **Kafka (Sarama)** — асинхронная обработка событий заказов.
- **gRPC + Protobuf** — интеграция с внешним сервисом заказов.
- **Easyp** — удобный Proto-менеджер для генерации, lint и backward-compatibility проверок.
- **Prometheus** — сбор метрик приложения.
- **Grafana** — визуализация метрик и мониторинг.
- **Zap** — структурированное логирование.
- **pprof** — профилирование производительности.
- **Docker / Docker Compose** — контейнеризация и локальный запуск окружения.

## 1. О проекте
`service-courier` — production-ready микросервис на Go для управления курьерами в системе доставки еды.
Он отвечает за жизненный цикл курьеров, автоматическое назначение курьера на заказ, снятие/завершение доставки, контроль «зависших» доставок,
интеграции с внешним сервисом заказов и сбор технических метрик.

Сервис спроектирован как компонент микросервисной архитектуры, ориентированный на масштабируемость, наблюдаемость и надежную работу под нагрузкой.

<img width="980" height="402" alt="image" src="https://github.com/user-attachments/assets/19ec93dd-fbdf-4ad0-b4c4-492729ec70f1" />


## 2. Ключевые цели проекта
В проекте применяются ключевые backend-практики:
- проектирование и поддержка микросервисной архитектуры;
- чистый и структурированный код (слои, интерфейсы, DI, SOLID);
- работа с PostgreSQL и миграциями;
- конкурентный код (goroutine, ticker, context, graceful shutdown);
- синхронная и асинхронная интеграция с внешними сервисами;
- observability (логирование, Prometheus, Grafana, pprof);
- тестирование (unit + integration), линтинг, CI.

## 3. Предметная область и роль сервиса
Сервис встраивается в экосистему food-delivery и покрывает следующие бизнес-задачи:
1. Управление курьерами: добавление, обновление, получение одного и списка.
2. Автоматическое назначение свободного курьера на заказ.
3. Снятие курьера с заказа и перевод в доступный статус.
4. Автоматическое освобождение занятых курьеров, если дедлайн доставки истёк.
5. Событийная обработка статусов заказов из Kafka.
6. Техническая эксплуатация: метрики, лимитирование, retry, профилирование.

## 4. Архитектура проекта
Проект организован по слоистой архитектуре с явным разделением ответственности:

- `cmd/app` — HTTP-приложение (API + фоновые воркеры + pprof);
- `cmd/consumer` — Kafka consumer + gRPC gateway к service-order;
- `internal/handler` — транспортный слой (HTTP/Kafka handlers);
- `internal/service` — бизнес-логика (usecase/services, фабрики);
- `internal/repository` — доступ к PostgreSQL;
- `internal/domain` — доменные модели и ошибки;
- `internal/gateway` — интеграции с внешними сервисами;
- `internal/worker` — фоновые процессы;
- `internal/middleware` — middleware/interceptor для логов, метрик, limiter, retry;
- `observability` — конфиги/адаптеры логирования и метрик;
- `deploy` — Dockerfile и compose для app/consumer/observability;
- `migrations` — SQL-миграции Goose.

## 5. Реализованный функционал по этапам

### 5.1 Базовый HTTP-сервис и graceful shutdown
- Поднят HTTP-сервер на `chi`.
- Реализованы health endpoints:
  - `GET /ping` → `{ "message": "pong" }`
  - `HEAD /healthcheck` → `204 No Content`
- Конфигурация читается из `.env`, порт можно переопределить через CLI-флаг `--port`.
- Корректная остановка через `signal.NotifyContext` и graceful shutdown.

### 5.2 PostgreSQL, миграции и CRUD курьеров
- Подключение к PostgreSQL через `pgxpool`.
- Миграции через `goose`.
- Реализован CRUD для сущности курьера:
  - `POST /courier`
  - `PUT /courier`
  - `GET /courier/{id}`
  - `GET /couriers`
- Добавлены валидации и маппинг доменных/транспортных ошибок в HTTP-ответы.

### 5.3 Слоистая архитектура и dependency injection
- Код разделён на handler/service/repository/domain.
- Интерфейсы определены на границах слоёв.
- Зависимости внедряются через конструкторы (`New...`).
- Сборка графа зависимостей выполняется в точках входа `main.go`.

### 5.4 Назначение/снятие курьера с заказа
- Добавлены endpoint’ы:
  - `POST /delivery/assign`
  - `POST /delivery/unassign`
- В `couriers` добавлено поле `transport_type` (`on_foot`, `scooter`, `car`).
- Добавлена таблица `delivery`.
- Реализована `DeliveryTimeFactory`:
  - `on_foot` → +30 мин;
  - `scooter` → +15 мин;
  - `car` → +5 мин.
- Операции назначения/снятия сделаны атомарно через транзакции.

### 5.5 Конкурентность и фоновые воркеры
- Фоновый воркер `delivery monitor` запускается в goroutine.
- По тикеру проверяет просроченные доставки и освобождает занятых курьеров одним SQL-обновлением.
- Назначение курьера выбирает доступного с минимальным количеством доставок (одним SQL-запросом).
- Управление жизненным циклом воркеров через `context.Context`.

### 5.6 Unit и интеграционные тесты
- Добавлены unit-тесты для service/handler/retry/limiter и др.
- Добавлены integration-тесты для сценариев с БД.
- В проекте есть команды для запуска тестов и coverage-отчёта.

### 5.7 Интеграция с внешним сервисом заказов
- Реализован gateway-слой для работы с сервисом заказов через gRPC.
- Подготовлен воркер polling-модели (оставлен в коде как часть пройденного этапа).

### 5.8 Асинхронная интеграция через Kafka
- Добавлен отдельный consumer-процесс (`cmd/consumer`).
- Consumer слушает Kafka topic изменений заказов.
- Реализована стратегия обработки событий по статусам заказа:
  - `created` → назначение курьера;
  - `deleted` → снятие курьера;
  - `completed` → завершение доставки.
- Для защиты от race/out-of-order событий используется проверка актуального статуса заказа через gateway.

### 5.9 Observability: метрики, логирование, мониторинг
- Middleware логирует HTTP-запросы (method, path, status, duration).
- Добавлен `/metrics` для Prometheus.
- Для consumer добавлены gRPC-метрики и метрики ретраев.
- Подготовлен docker-compose для локального запуска Prometheus + Grafana.

### 5.10 Rate Limiter и Retry
- Реализован собственный token-bucket limiter в HTTP middleware.
- При превышении лимита отдаётся `429 Too Many Requests`.
- В gRPC-вызовах к внешнему сервису добавлен retry с backoff + jitter через unary-интерцептор

### 5.11 CI/CD и качество кода
- Подготовлены линтерные/генерационные команды в `Makefile`.
- Проект ориентирован на подключение CI с запуском линтера и тестов.

### 5.12 Профилирование и оптимизация
- Поднят отдельный pprof-сервер с endpoint’ами `/debug/pprof/*`.
- В репозитории есть сохранённые профили (`profiles/*.pprof`, `trace.out`).
- Добавлена миграция с индексом для ускорения запросов по delivery.

## 6. API сервиса

### Health
- `GET /ping`
- `HEAD /healthcheck`

### Курьеры
- `POST /courier`
  ```json
  {
    "name": "Антон",
    "phone": "+37444111222",
    "status": "available",
    "transport_type": "on_foot"
  }
  ```
- `PUT /courier`
  ```json
  {
    "id": 1,
    "name": "Антон Иванов",
    "phone": "+37444111222",
    "status": "paused",
    "transport_type": "scooter"
  }
  ```
- `GET /courier/{id}`
- `GET /couriers`

### Доставки
- `POST /delivery/assign`
  ```json
  { "order_id": "8e6f9097-7c2e-4d84-ba28-0f3b5521a09c" }
  ```
- `POST /delivery/unassign`
  ```json
  { "order_id": "8e6f9097-7c2e-4d84-ba28-0f3b5521a09c" }
  ```

### Технические endpoints
- `GET /metrics`
- `GET /debug/pprof/*` (отдельный pprof-сервер)

## 7. Хранилище данных

### Таблица `couriers`
- `id`, `name`, `phone`, `status`, `transport_type`, `created_at`, `updated_at`.

### Таблица `delivery`
- `id`, `courier_id`, `order_id`, `assigned_at`, `deadline`.

### Индексы
- `idx_delivery_courier_id` на `delivery(courier_id)`.

## 8. Конфигурация

Ключевые переменные окружения (используются `cmd/app` и `cmd/consumer`):

### App (HTTP)
- `COURIER_LOCALHOST` (по умолчанию `localhost`)
- `COURIER_LOCALPORT` (по умолчанию `8080`)
- `TIME_CHECK` (период проверки дедлайнов)
- `TIME_POLL` (период polling-заказов, исторический этап)
- `REFILL` (интервал пополнения token bucket)
- `LIMIT` (ёмкость/лимит bucket)

### PostgreSQL
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_DB`
- `POSTGRES_HOST`
- `POSTGRES_LOCALPORT`
- `POSTGRES_PORT` (для docker-публикации порта)

### Consumer + Kafka + gRPC
- `CONSUMER_LOCALHOST`, `CONSUMER_LOCALPORT`
- `ORDER_HOST`, `ORDER_GRPC_PORT`
- `KAFKA_HOST`, `KAFKA_PORT`, `KAFKA_TOPIC`, `CONSUMER_GROUP`
- `COMMIT_INTERVAL`
- `MAX_ATTEMPTS`, `MULTIPLIER`, `JITTER`, `INIT_DELAY`, `MAX_DELAY`

### Observability
- `PROMETHEUS_PORT`, `PROMETHEUS_LOCALPORT`
- `GRAFANA_PORT`, `GRAFANA_LOCALPORT`

## 9. Локальный запуск

> Ниже — базовый сценарий локального запуска через существующие compose-файлы проекта.

### 1) Установите зависимости Go
```bash
go mod download
go mod tidy
```

### 2) Установите инструменты для работы с protobuf (из Makefile)
```bash
make bin-deps
```

Команда устанавливает в локальную папку `bin/`:
- `protoc-gen-go`
- `protoc-gen-go-grpc`
- `easyp`

Полезные команды из Makefile для protobuf:
```bash
make generate   # генерация protobuf/gRPC кода через easyp
make lint       # проверка proto-схем через easyp lint
make breaking   # проверка обратной совместимости proto-контрактов
```

### 3) Подготовьте инфраструктурную сеть
Проект ожидает внешнюю Docker-сеть `infrastructure_default`.

Если у вас уже есть отдельный репозиторий `infrastructure` — поднимите его:
```bash
make up-infra
```

Или создайте сеть вручную:
```bash
docker network create infrastructure_default
```

### 4) Создайте `.env` в нужных каталогах
Файлы compose читают `.env` из директории, где запускается команда.

Минимальный пример для `deploy/app/.env`:
```dotenv
# app
COURIER_LOCALHOST=0.0.0.0
COURIER_LOCALPORT=8080
COURIER_PORT=8082
TIME_CHECK=10s
TIME_POLL=5s
REFILL=1s
LIMIT=5

# postgres
POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypassword
POSTGRES_DB=test_db
POSTGRES_HOST=postgres
POSTGRES_LOCALPORT=5432
POSTGRES_PORT=5432
```

Пример для `deploy/consumer/.env` (дополнительно):
```dotenv
CONSUMER_LOCALHOST=0.0.0.0
CONSUMER_LOCALPORT=8083

ORDER_HOST=service-order
ORDER_GRPC_PORT=9000

KAFKA_HOST=kafka
KAFKA_PORT=9092
KAFKA_TOPIC=order.changed
CONSUMER_GROUP=service-courier
COMMIT_INTERVAL=1s

MAX_ATTEMPTS=3
MULTIPLIER=2
JITTER=0.1
INIT_DELAY=200ms
MAX_DELAY=5s

POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypassword
POSTGRES_DB=test_db
POSTGRES_HOST=postgres
POSTGRES_LOCALPORT=5432
```

Пример для `deploy/observability/.env`:
```dotenv
PROMETHEUS_PORT=9090
PROMETHEUS_LOCALPORT=9090
GRAFANA_PORT=3030
GRAFANA_LOCALPORT=3030
```

### 5) Запустите HTTP-приложение и БД
```bash
make up-app
```

### 6) Примените миграции
Используйте `goose` любым удобным способом (локально или в контейнере), указывая директорию `migrations/` и строку подключения к вашей БД.

### 7) (Опционально) Запустите consumer
```bash
make up-consumer
```

### 8) (Опционально) Запустите мониторинг
```bash
make up-observability
```
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3030`

### 9) Проверьте API
```bash
curl -i http://localhost:8082/ping
curl -I http://localhost:8082/healthcheck
```

### 10) Остановка
```bash
make down-consumer
make down-observability
make down-app
make down-infra
```

### 11) Запуск тестов
- Полный запуск тестов + race + coverage:
  ```bash
  make test
  ```
- Unit/integration вручную:
  ```bash
  go test ./...
  ```
  
---

### Автор
**Алексей Малков - [Ссылка на GitHub](https://github.com/shft1)**.
