# CLAUDE.md — go-blueprint

## Project Overview

Production-ready Go microservice blueprint. Demonstrates layered architecture with REST, gRPC, and event-driven consumers, full OpenTelemetry observability, and multiple data stores.

- **Module:** `github.com/stepanbukhtii/go-blueprint`
- **Go version:** 1.25.0
- **Entry point:** `main.go` → `cmd/root.go` (Cobra CLI)

---

## Essential Commands

```bash
# Build & run
go build -o app . && ./app server rest
./app server grpc
./app consumer kafka
./app consumer rabbitmq
./app consumer nats

# Testing
make test                    # go test -race ./...

# Linting
make lint                    # golangci-lint run

# Database
make migrate                 # apply migrations (goose up)
make migrate-down            # rollback (goose down)
make migrate-create          # new migration file
make model                   # setup DB + run migrations + regenerate SQLBoiler models

# Code generation
make generate                # go generate ./...
make docs                    # swag init (regenerates docs/swagger.yaml)
```

---

## Architecture

```
main.go
└── cmd/                     # Cobra CLI commands
    ├── server/rest.go        # HTTP server startup
    ├── server/grpc.go        # gRPC server startup
    └── consumer/{kafka,rabbitmq,nats}.go

internal/
├── app/                     # Wiring: DB, Redis, brokers, OTel, DI
│   ├── app.go               # App struct + infrastructure init
│   └── services.go          # samber/do dependency injection
├── config/                  # Typed env-var config
├── domain/                  # Interfaces + domain models (no deps)
├── repository/              # Data access
│   ├── postgres/            # SQLBoiler ORM implementations
│   │   ├── model/           # Auto-generated SQLBoiler models (do not edit)
│   │   └── convert/         # Domain ↔ model mappers
│   └── cached/              # Redis decorator over repository
├── service/                 # Business logic
│   ├── aggregator/          # Complex queries joining multiple domains
│   └── events/              # Event constants + factory functions
├── clients/                 # External service wrappers
│   ├── randomuser/          # HTTP client (randomuser.me)
│   └── user/                # gRPC client (remote user service)
└── transport/               # I/O adapters
    ├── http/                # Gin REST + Swagger + JWT middleware
    ├── grpc/                # gRPC server
    ├── kafka/               # franz-go consumer group
    ├── rabbitmq/            # amqp091 consumer
    └── nats/                # NATS subscriber

pkg/grpc/proto/              # Protobuf definitions + generated code
migrations/                  # Goose SQL migrations
docs/                        # Embedded Swagger YAML
deployments/                 # docker-compose (Grafana, Alloy, Prometheus, Loki, Tempo)
```

**Dependency direction:** `transport → service → repository → domain` (domain has no deps).

---

## Key Patterns

### Dependency Injection
Uses `github.com/samber/do/v2`. All services/repos registered in `internal/app/services.go`. Resolve with `do.MustInvoke[T](injector)`.

### Repository Pattern
`internal/repository/repo.go` defines the top-level `Repository` interface. PostgreSQL implementations live under `postgres/`; `cached/user_type.go` wraps the Postgres repo with Redis caching as a decorator.

### SQLBoiler ORM
Models in `internal/repository/postgres/model/` are **auto-generated** — run `make model` to regenerate. Never edit those files manually. Domain↔model conversion is in `postgres/convert/`.

### Event-Driven Flow
| Event | Producer | Consumer transport | Handler |
|---|---|---|---|
| `user.created` | `service/user.go` | Kafka | `transport/kafka/handlers/user_created.go` |
| `user.updated` | `service/user.go` | RabbitMQ | `transport/rabbitmq/handlers/user_updated.go` |
| `company.updated` | `service/company.go` | NATS | `transport/nats/handlers/company.go` |

Event constants and factory functions are in `internal/service/events/`.

### Aggregator
`service/aggregator/user.go` provides `UserAggregator` interface for fetching users with nested relationships (manager companies, user types) without polluting the core `UserService`.

### HTTP Transport
- Framework: Gin
- JWT middleware in `transport/http/middleware.go`
- Routes registered in `transport/http/handlers.go`
- Request/response DTOs in `transport/http/handlers/request/` and `response/`
- Swagger auto-generated; embed via `docs/embed.go`

### Observability (OpenTelemetry)
Initialized in `internal/app/app.go`. Exporters send to Alloy (OTLP gRPC):
- **Traces** — otelpgx (Postgres), gRPC, HTTP
- **Metrics** — periodic reader
- **Logs** — batch processor + stdout; bridged from `slog`

Local stack: `docker compose -f deployments/docker-compose.yaml up`
- Grafana: http://localhost:3000
- Prometheus: http://localhost:9090

---

## Domain Models

```go
// User — primary aggregate
type User struct {
    ID, Name, Username, Password, PublicName, Description string
    UserType       UserType
    Age, InitialAge int
    Rate, LastRate  decimal.Decimal
    Balance, LockBalance decimal.Decimal
    IsActive, ReadMessage bool
    LastLogin, CreatedAt, UpdatedAt time.Time
    ManagerCompanies []Company
}

// Company
type Company struct {
    ID, Name, OwnerID, ManagerID string
    IsActive bool
    LogoURL  *string
    Owner, Manager *User
    CreatedAt, UpdatedAt time.Time
}

// UserType — cached in Redis
type UserType struct {
    Code    string   // "DEFAULT" | "ADMIN"
    Name    string
    IsAdmin bool
}
```

---

## REST API (base: `/api/v1`)

| Method | Path | Auth |
|---|---|---|
| POST | `/auth/login` | — |
| GET/POST | `/users` | — |
| GET/PATCH/DELETE | `/users/:user_id` | — |
| GET/POST | `/user-types` | — |
| GET/PATCH/DELETE | `/user-types/:code` | — |
| GET/POST | `/companies` | JWT |
| GET/PATCH/DELETE | `/companies/:id` | JWT (admin for PATCH/DELETE) |
| GET | `/companies/owner` | JWT |
| POST | `/companies/multiple` | JWT |

Swagger UI available when `API_SWAGGER_ENABLE=true` → `GET /swagger/*`.

---

## gRPC API

Proto: `pkg/grpc/proto/user.proto`

```protobuf
service UserService {
  rpc One(OneRequest) returns (User);
}
```

---

## Configuration (`.env`)

All config via environment variables. See `.example.env` for all keys.
Key groups: `SERVICE_*`, `API_*`, `JWT_*`, `DB_*`, `REDIS_*`, `RABBITMQ_*`, `KAFKA_*`, `NATS_*`, `OTEL_*`, `RANDOM_USER_*`.

Config struct: `internal/config/config.go`.

---

## Database

- **Driver:** `pgx/v5` with connection pool (`pgxpool`)
- **ORM:** SQLBoiler (code-gen, not active-record)
- **Migrations:** Goose (`migrations/*.sql`)
- **Schema:** `user_type` (code PK, name, is_admin) + `users` (UUID PK, FK to user_type)

To add a table: write migration → `make model` to regenerate ORM models → add domain interface → implement repo + converter.

---

## Adding a New Feature (checklist)

1. Add domain model + interfaces to `internal/domain/`
2. Write SQLBoiler-based repo in `internal/repository/postgres/`
3. Add domain↔model converters in `postgres/convert/`
4. Register repo in `internal/repository/repo.go`
5. Implement service in `internal/service/`
6. Register service in `internal/app/services.go` (samber/do)
7. Add HTTP handler + request/response DTOs in `transport/http/handlers/`
8. Register route in `transport/http/handlers.go`
9. Run `make docs` to update Swagger

---

## Testing

Tests use `testify`. Run with race detector: `make test`.
Example: `internal/clients/randomuser/client_test.go`.
No mocks for DB by convention — integration tests preferred.

---

## Linting

Config: `.golangci.yml`. Active linters include `gosec`, `gocyclo`, `funlen`, `lll` (120 chars), `revive`. Test files are excluded from most rules.

Run: `make lint`
