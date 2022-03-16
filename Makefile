LDFLAGS=-ldflags="-s -w"

default: build

all: vet test build build-view

vet:
	go vet

test:
	go test ./... -v -timeout 10m

build:
	go build ${LDFLAGS} -o ./bin/webServer ./cmd/webServer/webServer.go
	go build ${LDFLAGS} -o ./bin/dataServer ./cmd/dataServer/dataServer.go
	go build ${LDFLAGS} -o ./bin/util ./cmd/util/util.go

build-web:
	cd web && yarn && yarn prod

run:
	cd bin && ./server

dev:
	cd bin && ./server -m development

dev-web:
	cd web && yarn && yarn dev