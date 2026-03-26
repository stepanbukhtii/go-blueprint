# Commands

```
go tool goose -dir ./migrations postgres "postgres://postgres:password@localhost:5432/postgres" up
```

```
go tool github.com/aarondl/sqlboiler/v4 -c ./sqlboiler.toml psql
```

```
go tool swag init -g ./internal/transport/http/server.go --parseDependency --parseInternal --outputTypes yaml --output docs
```

```
go tool swag fmt
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
- GRPC

- claude md