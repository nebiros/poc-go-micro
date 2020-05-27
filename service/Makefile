
GOPATH:=$(shell go env GOPATH)
MODIFY=Mgithub.com/micro/go-micro/api/proto/api.proto=github.com/micro/go-micro/v2/api/proto

.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/service/service.proto

.PHONY: build
build: proto
	go build -race -v -o service *.go

.PHONY: test
test:
	go test -race -v ./... -cover

.PHONY: docker
docker:
	docker build . -t service:latest

.PHONY: run
run:
	go run -race -v main.go
