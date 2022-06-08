.PHONY: test
test:
	@echo "\n 🛠️  Running unit tests..."
	go test ./...

.PHONY: build
build:
	@echo "\n 🔧  Building binary..."
	GOOS=darwin GOARCH=amd64 go build -a -o bin/ktop cmd/ktop.go
