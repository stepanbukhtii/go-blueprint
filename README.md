# Commands

## Migrations
```
go tool goose -dir ./migrations postgres "postgres://postgres:password@localhost:5432/postgres" up
```

```
go tool github.com/aarondl/sqlboiler/v4 -c ./sqlboiler.toml psql
```

## Swagger
```
go tool swag fmt
```

```
go tool swag init -g ./internal/transport/http/server.go --parseDependency --parseInternal --outputTypes yaml --output docs
```

## GRPC
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./pkg/grpc/proto/user.proto
```

## Stack

Api Router - Gin
Web Server - NGINX
Collector - Alloy
Log System - Loki
Metrics System - Prometheus
Tracing System - Tempo
Data Visualization - Grafana


## TODO

- claude md