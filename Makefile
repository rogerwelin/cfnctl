check: test lint

.PHONY: test
test:
	go test -cover -race -v ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: build
build:
	CGO_ENABLED=0 go build ./cmd/cfnctl

# ==================================================================================== #
#  QUALITY CONTROL
# ==================================================================================== #

.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Linting code...'
	golangci-lint run ./...
