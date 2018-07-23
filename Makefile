VERSION = $(shell git rev-parse --short HEAD)

test:
	go test -v -timeout=30s github.com/elvis-yan/geoip/pkg

install:
	go install  -ldflags "-X main.version=${VERSION}" github.com/elvis-yan/geoip/cmd/
	mv ${GOPATH}/bin/cmd ${GOPATH}/bin/geoip

build:
	go build -ldflags "-X main.version=${VERSION}" -o geoip ./cmd/main.go
