export PATH := $(GOPATH)/bin:$(PATH)
BIN=./build

ifeq ($(OS),Windows_NT)
	EXT=.exe
else
	EXT=
endif

all: subserver

subserver:
	go build -o $(BIN)/subserver$(EXT)  ./
