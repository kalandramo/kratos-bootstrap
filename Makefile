.PHONY: api init plugin cli

# generate protobuf api go code
api:
	cd api && \
	buf generate

# initialize develop environment
init: plugin cli

# install protoc plugin
plugin:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v3@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v3@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/menta2k/protoc-gen-redact/v3@latest
	go install github.com/go-kratos/protoc-gen-typescript-http@latest

# install cli tools
cli:
	go install github.com/go-kratos/kratos/cmd/kratos/v3@latest
	go install github.com/google/gnostic@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install entgo.io/ent/cmd/ent@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest