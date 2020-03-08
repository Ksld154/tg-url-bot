
GOPATH = $(shell go env GOPATH)
export PATH := $(PATH):$(GOPATH)/bin

all: build

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/echobot cmd/echobot/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/webhookServer cmd/webhookServer/main.go

clean: 
	rm -rf ./bin 