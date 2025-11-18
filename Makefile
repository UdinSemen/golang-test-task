MOCKERY_VERSION=v3.5.5
LOCAL_BIN:=$(CURDIR)/bin
MOCKERY_BIN := $(LOCAL_BIN)/mockery
MOCKERY_CONFIG := .mockery.yml

.PHONY: swagger_init
swagger_init:
	swag fmt
	swag init -g cmd/app/main.go

.PHONY: migrations_create
migrations_create:
	migrate create -ext sql -dir ./migrations -seq ${NAME}

.PHONY: tests
tests:
	go test ./...

.PHONY: install-mockery
install-mockery: $(MOCKERY_BIN)
$(MOCKERY_BIN):
	@echo "🛠️ Installing mockery $(MOCKERY_VERSION)..."
	@GOBIN=$(LOCAL_BIN) go install github.com/vektra/mockery/v3@$(MOCKERY_VERSION)

.PHONY: generate-mocks
generate-mocks: install-mockery
	@echo "🔁 Generating mocks..."
	@$(MOCKERY_BIN) --config=$(MOCKERY_CONFIG)
