include .env

dbDSN := "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)"
migration := go tool goose -dir ./migrations -allow-missing postgres $(dbDSN)

model:
	@docker run --rm -d -p 5433:5432 --name migration-db \
		-e POSTGRES_PASSWORD="migration" \
		-e POSTGRES_DB="$(DB_NAME)" \
		-v ./deployments/local/init-postgres.sql:/docker-entrypoint-initdb.d/init.sql \
		postgres:18
	@sleep 5
	@go tool goose -dir ./migrations postgres $(dbDSN) up
	@sed 's/= \"dbname\"/= \"$(DB_NAME)\"/g; s/= \"user\"/= \"$(DB_USER)\"/g; s/= \"pass\"/= \"$(DB_PASSWORD)\"/g; s/= 5432/= $(DB_PORT_TEMP)/g; s/= \"schema\"/= \"$(DB_SCHEMA)\"/g' ./sqlboiler.toml > ./temp_boiler_config.toml
	@go tool github.com/aarondl/sqlboiler/v4 -c ./temp_boiler_config.toml psql
	@rm ./temp_boiler_config.toml
	@docker stop migration-db

migrate-create:
	$(migration) create $(name) sql

migrate:
	$(migration) up

migrate-down:
	$(migration) down

lint:
	golangci-lint run

generate:
	go generate -v ./...

test: generate
	go test -v -race ./...

docs:
	@rm -rf apispec/openapi/v1/swaggo
	@go tool swag fmt

	@go tool swag init -g main.go --parseDependency --parseInternal --outputTypes yaml --output docs
