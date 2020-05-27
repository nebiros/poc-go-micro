
GOPATH:=$(shell go env GOPATH)
MODIFY=Mgithub.com/micro/go-micro/api/proto/api.proto=github.com/micro/go-micro/v2/api/proto

.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/client/client.proto

.PHONY: build
build: proto
	go build -race -v -o client *.go

.PHONY: test
test:
	go test -race -v ./... -cover

.PHONY: docker
docker:
	docker build . -t client:latest

.PHONY: run
run:
	go run -race -v main.go
