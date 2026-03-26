# CLAUDE.md — go-blueprint

## Project

- **Module:** `github.com/stepanbukhtii/go-blueprint`
- **Go version:** 1.26.3
- **Entry point:** `main.go` → `cmd/root.go` (Cobra CLI)

---

## Commands

```bash
# Build & run
go build -o app . && ./app server rest
./app server grpc
./app consumer kafka
./app consumer rabbitmq
./app consumer nats

# Testing
make test                    # go generate + go test -race ./...

# Linting
make lint                    # golangci-lint run

# Database
make migrate                 # goose up
make migrate-down            # goose down
make migrate-create name=<migration_name>
make model                   # spin up temp DB + goose up + bobgen-psql

# Code generation
make generate                # go generate ./...
make docs                    # swag init → docs/swagger.yaml
```

---

## Folder Structure

```
main.go
└── cmd/
    ├── server/rest.go         # HTTP server startup
    ├── server/grpc.go         # gRPC server startup
    └── consumer/{broker}.go   # one file per message broker

internal/
├── app/
│   ├── app.go                 # App struct + infrastructure init (DB, Redis, brokers, OTel)
│   └── services.go            # samber/do DI wiring
├── config/                    # Typed env-var config struct
├── domain/                    # Interfaces + domain models (no external deps)
├── repository/
│   ├── postgres/
│   │   ├── models/            # Auto-generated bob models — do not edit
│   │   ├── dberrors/          # Auto-generated DB constraint error types — do not edit
│   │   ├── dbinfo/            # Auto-generated table/column metadata — do not edit
│   │   ├── factory/           # Auto-generated test factories — do not edit
│   │   └── convert/           # Domain ↔ model converters, one file per entity
│   └── cached/                # Redis decorator over postgres repositories
├── service/
│   ├── aggregator/            # Cross-domain queries joining multiple repos
│   └── events/                # Event name constants + message factory functions
├── clients/                   # External HTTP/gRPC service wrappers
└── transport/
    ├── http/
    │   ├── handlers/
    │   │   ├── request/       # Request DTOs with ToDomain() methods
    │   │   ├── response/      # Response DTOs with NewX(domain) constructors
    │   │   └── ws/            # WebSocket handlers
    │   ├── handlers.go        # Route registration
    │   └── middleware.go      # JWT middleware
    ├── grpc/
    │   └── handlers/
    │       └── response/      # gRPC response mappers
    ├── kafka/handlers/        # one file per event type
    ├── rabbitmq/handlers/
    └── nats/handlers/

pkg/grpc/proto/                # Protobuf definitions + generated code
migrations/                    # Goose SQL files, sequential numbering
docs/                          # Embedded Swagger YAML
deployments/                   # docker-compose, local init SQL
```

**Dependency direction:** `transport → service → repository → domain`  
`domain` has zero external dependencies.

---

## Code Conventions

### Domain Layer (`internal/domain/`)

One file per entity. Each file contains:
- Sentinel errors defined at package level using `errx.Wrap`:
  ```go
  var ErrFooNotFound = errx.Wrap(api.ErrNotFound, "foo not found")
  ```
- Domain struct (plain Go struct, no ORM tags)
- `FooRepository` interface — data access contract
- `CreateFooInput` / `UpdateFooInput` structs — service-layer inputs
- `FooService` interface — business logic contract
- `FooAggregator` interface (if cross-domain enrichment is needed)

Lookup/enum entities use a string code as primary key and expose constants at package level.

### Repository Layer (`internal/repository/postgres/`)

- Unexported struct with `exec bob.Executor` field
- Exported constructor returns the domain interface:
  ```go
  func NewFoo(exec bob.Executor) domain.FooRepository { return &foo{exec: exec} }
  ```
- `repository.Repository` (in `repo.go`) aggregates all domain repos and exposes `RunInTransaction`
- Standard methods: `Add`, `Update`, `Save`, `Find`, `FindAll`, `FindAllPaginate`, `Exists`, `Remove`
- `Save` = check `Exists` → `Add` or `Update`
- `FindAllPaginate` runs `Count` first, then applies `sm.Limit` + `sm.Offset`
- Translate `sql.ErrNoRows` to the domain sentinel error in `Find`

### Cached Decorator (`internal/repository/cached/`)

- Unexported struct with `cache` and `repo` fields
- Constructor: `NewFooRepository(cache cache.MapCache[domain.Foo], repo domain.FooRepository) domain.FooRepository`
- Cache-aside pattern: read from cache, fall back to repo, refresh cache on miss

### Converters (`internal/repository/postgres/convert/`)

One file per entity. Package-level zero-value converter variable used as a namespace:
```go
var Foo foo
type foo struct{}

func (foo) Domain(m *models.Foo) domain.Foo { ... }
func (c foo) DomainSlice(m models.FooSlice) []domain.Foo { return lo.Map(m, ...) }
func (foo) Setter(f *domain.Foo) *models.FooSetter { ... }
```

Null handling in setters:
- Non-null: `omit.From(val)`
- Nullable from pointer: `omitnull.FromPtr(ptr)`
- Nullable with condition: `omitnull.FromNull(null.FromCond(val, condition))`

### Service Layer (`internal/service/`)

- Constructor receives `do.Injector`, returns `(domain.XService, error)`:
  ```go
  func NewFooService(injector do.Injector) (domain.FooService, error) {
      return &fooService{
          repo: do.MustInvoke[repository.Repository](injector),
      }, nil
  }
  ```
- Structured logging with `slog` using dotted key paths (`"entity.id"`, `"event.name"`):
  ```go
  slog.With(slog.String("foo.id", foo.ID)).InfoContext(ctx, "foo created")
  ```
- Publish events after mutating state; log the publish separately

### Aggregator (`internal/service/aggregator/`)

Use for cross-domain enrichment. Pattern:
```go
func (a *fooAggregator) Get(ctx context.Context, id string) (domain.Foo, error) {
    foo, err := a.repo.Foo().Find(ctx, ...)
    if err != nil { return domain.Foo{}, err }
    if err := a.setOneRelations(ctx, &foo); err != nil { return domain.Foo{}, err }
    return foo, nil
}
```
Separate `setOneRelations`, `setManyRelations`, and per-relation helpers (`setFooBars`).  
For bulk enrichment, fetch the related slice once, build a `lo.KeyBy` map, then iterate.

### Events (`internal/service/events/`)

- `events.go` — string constants and payload structs:
  ```go
  const FooCreatedEvent = "foo.created"
  type EventFooCreatedData struct { FooID string `json:"foo_id"` }
  ```
- `converter.go` — factory functions `NewEventFooCreatedData(foo domain.Foo) EventFooCreatedData`

### HTTP Transport (`internal/transport/http/`)

- Handler struct holds only the domain service interface:
  ```go
  type Foo struct { service domain.FooService }
  func NewFoo(service domain.FooService) *Foo { return &Foo{service: service} }
  ```
- Route groups registered as `registerFooHandlers(r *gin.RouterGroup)` methods on `*Server`
- JWT middleware applied at the group level: `r.Group("/foos", s.jwtMiddleware.Auth)`
- Role checks added per-route: `group.PATCH("", h.Update, s.jwtMiddleware.AuthRole(domain.RoleAdmin))`
- Nested resource groups:
  ```go
  foos := r.Group("/foos")
  fooID := foos.Group("/:foo_id")
  ```
- Use `api.ParseRequest`, `api.RespondData`, `api.RespondDataPages`, `api.RespondOK`, `api.ServeError` from `easy-tools`
- Swagger annotations on every handler method

**Request DTOs** (`handlers/request/`):
- Embed URI struct for path params: `FooURI struct { FooID string \`uri:"foo_id"\` }`
- `ToDomain()` method converts DTO → domain input struct
- Validation via `binding` tags

**Response DTOs** (`handlers/response/`):
- Constructor `NewFoo(f domain.Foo) Foo` and `NewFoos(fs []domain.Foo) []Foo` (using `lo.Map`)
- Never reuse domain structs as response types

### gRPC Transport (`internal/transport/grpc/`)

- Handler struct embeds `proto.UnimplementedXServiceServer`
- Response mapped via `response.NewX(domain)` in `handlers/response/`
- Run `make generate` after changing `.proto` files

### DI (`internal/app/services.go`)

- `do.ProvideValue` for infrastructure (db, redis, brokers, config)
- `do.Provide` for services and aggregators (constructor-based)
- `do.MustInvoke[T]` to resolve and assign to `Services` struct
- Register aggregators before services (services may depend on them)

---

## bob ORM

Config file: `bobgen.yaml`.  
Auto-generated output directories: `models/`, `dberrors/`, `dbinfo/`, `factory/`. Never edit these manually.

Query patterns:
```go
// Insert
models.Foos.Insert(setter).One(ctx, exec)

// Update
models.Foos.Update(models.UpdateWhere.Foos.ID.EQ(id), setter.UpdateMod()).One(ctx, exec)

// Delete
models.Foos.Delete(models.DeleteWhere.Foos.ID.EQ(id)).Exec(ctx, exec)

// Select with where
models.Foos.Query(mods...).One(ctx, exec)
models.Foos.Query(sm.Where(models.Foos.Columns.ID.EQ(psql.Arg(id)))).One(ctx, exec)

// Count
models.Foos.Query(mods...).Count(ctx, exec)
models.Foos.Query(mods...).Exists(ctx, exec)
```

Import aliases: `sm` (select mods), `um` (update mods), `dm` (delete mods); args via `psql.Arg(val)`.

---

## Database

- **Driver:** `pgx/v5` via `pgxpool`; `stdlib` adapter for `*sql.DB`
- **Migrations:** Goose, `migrations/*.sql`, sequential numbering
- Primary keys: UUID strings for domain entities; natural string codes for lookup tables
- UUID columns mapped to `string` in bob (configured in `bobgen.yaml` replacements)

To add a table: write migration → `make model` → implement domain interface + repo + converter.

---

## Configuration

All config via environment variables; no config files at runtime. See `.example.env` for keys.  
Parsed into a typed struct in `internal/config/config.go` using `caarlos0/env`.  
Group env vars by prefix per concern (`DB_*`, `REDIS_*`, `KAFKA_*`, `JWT_*`, etc.).  
Project-specific config blocks are added as nested structs in `Config`.

---

## Testing

- Framework: `testify`
- Race detector: `make test` runs `go test -race ./...`
- No DB mocks — prefer integration tests against a real database
- Test files co-located with the package (`foo_test.go`)
- `factory/` package provides bob-generated test fixtures

---

## Linting

Config: `.golangci.yml`. Key linters: `gosec`, `gocyclo`, `funlen`, `lll` (120 chars), `revive`.  
Test files excluded from most rules. Run: `make lint`