APP=$(shell basename $(shell git remote get-url origin))
REGISTRY=ryakovlev
VERSION=$(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)

# change TARGETOS for different OS (linux, ios, windows)
TARGETOS=linux
# change TARGETARCH for different architectures (amd64, arm64)
TARGETARCH=amd64

lint:
	golangci-lint run

test:
	go test -v

get:
	go get

format:
	gofmt -s -w ./
	
build: format
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o kbot -ldflags "-X github.com/yakovlevruslan/kbot/cmd.appVersion=${VERSION}"

image:
	docker build . -t ${REGISTRY}/${APP}:${VERSION}-${TARGETARCH}-${TARGETOS}

push:
	docker push ${REGISTRY}/${APP}:${VERSION}-${TARGETARCH}-${TARGETOS}

clean:
	rm -rf kbot
	docker rmi $(shell docker images --filter=reference="ryakovlev*/*:v1.0.3*" -q) -f