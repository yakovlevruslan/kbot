VERSION=$(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
TARGETOS=linux

lint:
	golangci-lint run

test:
	go test -v

go:
	go get
	
format:
	gofmt -s -w ./
build: format
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${shell dpkg --print-architecture} go build -v -o kbot -ldflags "-X github.com/yakovlevruslan/kbot/cmd.appVersion=${VERSION}"

clean:
	rm -rf kbot