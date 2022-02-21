LDFLAGS=-ldflags="-s -w"

default: build

all: vet test build build-view

vet:
	go vet

test:
	go test ./... -v -timeout 10m

pprof-profile:
	go tool pprof -http=:41001 http://localhost:42073/debug/pprof/profile

pprof-heap:
	go tool pprof -http=:41002 http://localhost:42073/debug/pprof/heap

build:
	go build ${LDFLAGS} -o ./bin/server ./cmd/server/server.go
	go build ${LDFLAGS} -o ./bin/util ./cmd/util/util.go
	go build ${LDFLAGS} -o ./bin/resizer ./cmd/resizer/resizer.go


build-web:
	cd web && yarn && yarn prod

run:
	cd bin && ./server

dev:
	cd bin && ./server -m development

dev-web:
	cd web && yarn && yarn dev

.EXPORT_ALL_VARIABLES:
MALLOC_ARENA_MAX=2