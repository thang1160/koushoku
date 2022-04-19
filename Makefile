ARCHITECTURES=386 amd64
LDFLAGS=-ldflags="-s -w"

default: build

all: vet test build build-view

vet:
	go vet

test:
	go test ./... -v -timeout 10m

build:
	$(foreach GOARCH,$(ARCHITECTURES),\
		$(shell export GOARCH=$(GOARCH))\
		$(shell go build $(LDFLAGS) -o ./bin/$(GOARCH)/webServer ./cmd/webServer/webServer.go)\
		$(shell go build $(LDFLAGS) -o ./bin/$(GOARCH)/dataServer ./cmd/dataServer/dataServer.go)\
		$(shell go build $(LDFLAGS) -o ./bin/$(GOARCH)/util ./cmd/util/util.go)\
	)\

build-web:
	cd web && yarn && yarn prod

run:
	cd bin && ./server

dev:
	cd bin && ./server -m development

dev-web:
	cd web && yarn && yarn dev