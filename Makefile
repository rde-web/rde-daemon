APP=rde-daemon
VERSION=$(shell git branch --show-current)

ifeq ($(VERSION),master)
VERSION=latest
endif

run:
	@go run ./cmd/rde-daemon/main.go

build:
	@docker build -f build/Dockerfile -t $(APP):$(VERSION) .

run.docker:
	@docker run $(APP):$(VERSION)

.PHONY: build run run.docker