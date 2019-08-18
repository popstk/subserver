export PATH := $(GOPATH)/bin:$(PATH)
BIN=./bin

ifeq ($(OS),Windows_NT)
	EXT=.exe
else
	EXT=
endif

all: subserver

proto:
	protoc -I/usr/local/include -I. \
      -I$(GOPATH)/src \
      -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
      --go_out=plugins=grpc:. \
      proto/query/query.proto
	protoc -I/usr/local/include -I. \
      -I$(GOPATH)/src \
      -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
      --grpc-gateway_out=logtostderr=true:. \
      proto/query/query.proto

subserver:
	go build -o $(BIN)/subserver$(EXT)  ./cmd/subserver


.PHONY: proto


