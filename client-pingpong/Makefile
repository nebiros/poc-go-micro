
GOPATH:=$(shell go env GOPATH)
MODIFY=Mgithub.com/micro/go-micro/api/proto/api.proto=github.com/micro/go-micro/v2/api/proto

.PHONY: build
build:
	go build -race -v -o client-pingpong *.go

.PHONY: test
test:
	go test -race -v ./... -cover

.PHONY: docker
docker:
	docker build . -t client-pingpong:latest

.PHONY: run
run:
	go run -race -v main.go
