check: test lint vet

test:
	go test -cover -race -v ./...

lint:
	golint ./...

vet:
	go vet ./...

build:
	go build -o cfnctl cmd/cfnctl/main.go
