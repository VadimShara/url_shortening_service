.PHONY: all generate

all: generate

generate:
	@echo "Generating gRPC code..."
	protoc -I api/proto api/proto/service.proto \
		--go_out=api/gen/go/ --go_opt=paths=source_relative \
		--go-grpc_out=api/gen/go/ --go-grpc_opt=paths=source_relative

gen: generate