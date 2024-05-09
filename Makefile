include .env
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up: ### Run docker compose
	docker compose up --build -d postgres && docker compose logs -f
.PHONY: compose-up

compose-up-integration-test: ### Run docker compose with integration test
	docker compose up --build --abort-on-container-exit --exit-code-from integration
.PHONY: compose-up-integration-test

compose-down: ### Down docker compose
	docker compose down --remove-orphans
.PHONY: compose-down

swag-v1: ### swag init
	swag init -g internal/controller/http/v1/router.go
.PHONY: swag-v1

run: ### run app
	go mod tidy && go mod download && \
	DISABLE_SWAGGER_HTTP_HANDLER='' GIN_MODE=debug CGO_ENABLED=1 go run ./cmd/app
.PHONY: run

docker-rm-volume: ### remove docker volume
	docker volume rm go-clean-template_pg-data
.PHONY: docker-rm-volume

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-hadolint: ### check by hadolint linter
	git ls-files --exclude='Dockerfile*' --ignored | xargs hadolint
.PHONY: linter-hadolint

linter-dotenv: ### check by dotenv linter
	dotenv-linter
.PHONY: linter-dotenv

test: ### run test
	go test -v -cover -race ./internal/...
.PHONY: test

integration-test: ### run integration-test
	go clean -testcache && go test -v ./integration-test/...
.PHONY: integration-test

mock: ### run mockgen
	mockgen -source ./internal/usecase/ciraconfigs/interfaces.go -package ciraconfigs_test > ./internal/usecase/ciraconfigs/mocks_test.go
	mockgen -source ./internal/usecase/devices/interfaces.go -package devices_test > ./internal/usecase/devices/mocks_test.go
	mockgen -source ./internal/usecase/domains/interfaces.go -package domains_test > ./internal/usecase/domains/mocks_test.go
	mockgen -source ./internal/usecase/ieee8021xconfigs/interfaces.go -package ieee8021xconfigs_test > ./internal/usecase/ieee8021xconfigs/mocks_test.go
	mockgen -source ./internal/usecase/profiles/interfaces.go -package profiles_test > ./internal/usecase/profiles/mocks_test.go
	mockgen -source ./internal/usecase/wificonfigs/interfaces.go -package wificonfigs_test > ./internal/usecase/wificonfigs/mocks_test.go

.PHONY: mock

migrate-create:  ### create new migration
	migrate create -ext sql -dir migrations 'migrate_name'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(DB_URL)?sslmode=disable' up
.PHONY: migrate-up

bin-deps:
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	GOBIN=$(LOCAL_BIN) go install go.uber.org/mock/mockgen@latest
