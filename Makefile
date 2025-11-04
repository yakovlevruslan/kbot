VERSION=$(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)

format:
	gofmt -s -w ./
build:
	go build -v -o kbot -ldflags "-X github.com/yakovlevruslan/kbot/cmd.appVersion=${VERSION}"