BUILD_DIR = builds
BINARY_NAME = gobot-pir
MODULE = github.com/soerenschneider/gobot-pir
CHECKSUM_FILE = $(BUILD_DIR)/checksum.sha256
SIGNATURE_KEYFILE = ~/.signify/github.sec

tests:
	go test ./...

clean:
	rm -rf ./$(BUILD_DIR)

build: version-info
	CGO_ENABLED=0 go build -ldflags="-X '$(MODULE)/internal.BuildTime=${BUILD_TIME}' -X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BINARY_NAME) cmd/main.go

release: clean version-info cross-build
	sha256sum $(BUILD_DIR)/$(BINARY_NAME)-* > $(CHECKSUM_FILE)
	gh-upload-assets -o soerenschneider -r gobot-pir -f ~/.gh-token builds

signed-release: release
	pass keys/signify/github | signify -S -s $(SIGNATURE_KEYFILE) -m $(CHECKSUM_FILE)

cross-build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0       go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64    cmd/main.go
	GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=0 go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-armv5     cmd/main.go
	GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-armv6     cmd/main.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0       go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-aarch64   cmd/main.go

version-info:
	$(eval VERSION := $(shell git describe --tags --abbrev=0 || echo "dev"))
	$(eval COMMIT_HASH := $(shell git rev-parse HEAD))

coverage:
	go test ./... -covermode=count -coverprofile=coverage.out
	go tool cover -html=coverage.out -o=coverage.html
	go tool cover -func=coverage.out -o=coverage.out
