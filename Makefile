PKG=gitlab.apulis.com.cn/hjl/blank-web-app
IMAGE?=apulistech/blankWebApp-backend
GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS?="-X ${PKG}/cmd.gitCommit=${GIT_COMMIT} -X ${PKG}/cmd.buildDate=${BUILD_DATE}"
GO111MODULE=on
GOPROXY=https://goproxy.io
GOPATH=$(shell go env GOPATH)

.EXPORT_ALL_VARIABLES:

.PHONY: build
build: generate
	mkdir -p bin
	GOOS=linux go build -ldflags ${LDFLAGS} -buildmode=pie -o bin/blankWebApp main.go

 # Regenerates OPA data from rego files
HAVE_GO_BINDATA := $(shell command -v go-bindata 2> /dev/null)
.PHONY: generate
generate:
ifndef HAVE_GO_BINDATA
	@echo "requires 'go-bindata' (go get -u github.com/kevinburke/go-bindata/go-bindata)"
	@exit 1 # fail
else
	go generate ./...
endif


.PHONY: verify
verify:
	./hack/verify-all

.PHONY: test
test:
	go test -v -race ./pkg/...

.PHONY: image
image:
	docker build -t $(IMAGE):${GIT_BRANCH}-${GIT_COMMIT} -f Dockerfile .
