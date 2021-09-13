BINARY_NAME = gobot-motion

build:
	go build -o $(BINARY_NAME) cmd/main.go

build-raspberry:
	GOARM=6 GOARCH=arm GOOS=linux go build -o $(BINARY_NAME)_armv6 cmd/main.go

unittest:
	go test ./...
