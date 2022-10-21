GRPC_DEPS:=\
	google.golang.org/protobuf/cmd/protoc-gen-go\
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
run:
	@go run ./cmd/daemon/main.go

docker-generate:
	@docker run --rm \
		-v$(CURDIR):/generate \
		dafaque/go-grpc-generator:0.0.4-alpine3.16
	@go mod tidy

generate:
	@protoc\
		--proto_path=api/rde-daemon-api\
		--go_out=internal\
		--go_opt=paths=import\
		--go-grpc_out=internal\
		--go-grpc_opt=paths=import\
		api/rde-daemon-api/*.proto

dependency-install:
	@go get $(GRPC_DEPS)
	@go install $(GRPC_DEPS)

modules:
	@git submodule update --init --remote api/