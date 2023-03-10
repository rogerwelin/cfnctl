check: test lint vet

.PHONY: test
test:
	go test -cover -race -v ./...

.PHONY: lint
lint:
	golint ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build:
	go build -o cfnctl cmd/cfnctl/main.go

