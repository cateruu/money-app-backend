include .envrc

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

.PHONY: run/api
run/api:
	@go run ./cmd/api -dsn=${DB_DSN}

.PHONY: psql
psql:
	@psql ${DB_DSN}

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build: build the cmd/api application
.PHONY: build/api
build/api: 
	@echo 'Building cmd/api...'
	go build -ldflags="-s -w" -o=./bin/api ./cmd/api