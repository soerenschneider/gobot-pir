build-raspberry:
	GOARM=6 GOARCH=arm GOOS=linux go build -o gobot-motion cmd/main.go

unittest:
	go test ./...
