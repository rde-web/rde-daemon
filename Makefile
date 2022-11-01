run:
	@go run ./cmd/rde-daemon/main.go

build:
	@docker build -f build/Dockerfile -t rde-daemon:local .

run.docker:
	@docker run rde-daemon:local

.PHONY: build run run.docker