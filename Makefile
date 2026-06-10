.PHONY: api init plugin cli tidy gen update-all

# generate protobuf api go code
api:
	cd api && buf generate

# generate all code
gen: api

# sync all workspace modules
tidy:
	@echo "==> syncing workspace..."
	go work sync
	@echo "==> tidying all modules..."
	@for mod in $$(find . -maxdepth 2 -name "go.mod" -not -path "./.claude/*"); do \
		dir=$$(dirname "$$mod"); \
		echo "==> tidy: $$dir"; \
		(cd "$$dir" && go mod tidy); \
	done
	@echo "==> done"

# update all dependencies to latest local versions
# usage: make update-all
update-all:
	@echo "==> updating all local module dependencies..."
	@# update api version in go.work
	@echo "==> updating api in workspace..."
	go work use ./api
	@# update bootstrap version in go.work
	@echo "==> updating bootstrap in workspace..."
	go work use ./bootstrap
	@# update server version in go.work
	@echo "==> updating server in workspace..."
	go work use ./server
	@# update config version in go.work
	@echo "==> updating config in workspace..."
	go work use ./config
	@# update example version in go.work
	@echo "==> updating example in workspace..."
	go work use ./example
	@# sync workspace
	@echo "==> syncing workspace..."
	go work sync
	@echo "==> done"

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
