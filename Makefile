BINARY_NAME = gobot-motion-detection
MODULE = gobot-motion-detection

build: version-info
	go build -ldflags="-X '$(MODULE)/internal.BuildTime=${BUILD_TIME}' -X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BINARY_NAME) cmd/main.go

build-raspberry: version-info
	GOARM=6 GOARCH=arm GOOS=linux go build -o $(BINARY_NAME)-armv6 cmd/main.go

version-info:
	$(eval VERSION := $(shell git describe --tags || echo "dev"))
	$(eval BUILD_TIME := $(shell date +"%Y-%m-%dT%H:%M:%SZ"))
	$(eval COMMIT_HASH := $(shell git rev-parse --short HEAD))

coverage:
	go test ./... -covermode=count -coverprofile=coverage.out
	go tool cover -html=coverage.out -o=coverage.html
	go tool cover -func=coverage.out -o=coverage.out

unittest:
	go test ./...
