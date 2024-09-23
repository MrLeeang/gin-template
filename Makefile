GOPATH:=$(shell go env GOPATH)

.PHONY: init
init:
	@go get -u google.golang.org/protobuf/proto
	@go install github.com/golang/protobuf/protoc-gen-go@latest
	@go install github.com/asim/go-micro/cmd/protoc-gen-micro/v4@latest

.PHONY: proto
proto:
	@protoc --proto_path=./service --micro_out=./service --go_out=:./service service/proto/service.proto
	@protoc --proto_path=./service --micro_out=./service --go_out=:./service service/proto/mail.proto
	@protoc --proto_path=./service --micro_out=./service --go_out=:./service service/proto/sms.proto

.PHONY: build-api
build-api:
	@go build -o api cmd/appv1/main.go

.PHONY: build-service
build-service:
	@go build -o service cmd/service/main.go

.PHONY: test
test:
	@go test -v ./... -cover

.PHONY: docker-api
docker-api:
	@docker build -f Dockerfile.api -t api:latest .

.PHONY: docker-service
docker-service:
	@docker build -f Dockerfile.service -t service:latest .
