.PHONY: go_all
go_all: go_clean go_tidy go_fmt go_vuln go_lint_clean_cache go_lint

.PHONY: go_tidy
go_tidy:
	@go mod tidy

.PHONY: go_fmt
go_fmt:
	@go fmt ./...

.PHONY: go_vuln
go_vuln:
	@go mod tidy
	@govulncheck ./...

.PHONY: go_lint_clean_cache
go_lint_clean_cache:
	@golangci-lint cache clean

.PHONY: go_lint
go_lint:
	@go mod tidy
	@golangci-lint run

.PHONY: go_lint_v
go_lint_v:
	@go mod tidy
	@golangci-lint run -v > /dev/null

.PHONY: go_lint_x
go_lint_x:
	@go mod tidy
	@golangci-lint run --no-config --enable-all -D=gofumpt,gci,varnamelen,exhaustivestruct,structcheck,nosnakecase,golint,maligned,scopelint,ifshort,varcheck,interfacer,deadcode,depguard

.PHONY: go_clean
go_clean:
	@go mod tidy
	@go clean -testcache -fuzzcache -cache

.PHONY: go_test_example_x
go_test_example_x:
	@go mod tidy
	@go test -timeout 90s -run ^ExampleX$

.PHONY: go_test_x
go_test_x:
	@go mod tidy
	@go clean -testcache -fuzzcache
	@go test -timeout 90s -run ^TestX$

.PHONY: go_test
go_test:
	@go mod tidy
	@go clean -testcache -fuzzcache
	@go test ./test

.PHONY: go_test_race
go_test_race:
	@go mod tidy
	@go clean -testcache -fuzzcache
	@go test -race ./test

.PHONY: go_bench_x
go_bench_x:
	@go mod tidy
	@go clean -testcache -fuzzcache
	@go test -benchmem -bench .

.PHONY: generate_protoc
generate_protoc:
	@protoc -I . -I /home/_/gh/googleapis/googleapis -I /home/_/gh/bufbuild/protoc-gen-validate \
		--go_out . \
		--go_opt paths=source_relative \
		\
 		--go-grpc_out . \
		--go-grpc_opt paths=source_relative \
		\
 		--grpc-gateway_out . \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		\
 		--validate_out . \
		--validate_opt paths=source_relative \
		--validate_opt lang=go \
		\
 		--openapiv2_out proto/ \
		--openapiv2_opt logtostderr=true \
		--openapiv2_opt enums_as_ints=false \
		--openapiv2_opt omit_enum_default_value=true \
		--openapiv2_opt allow_merge=true \
		--openapiv2_opt generate_unbound_methods=true \
		--openapiv2_opt repeated_path_param_separator=ssv \
		\
		proto/*/*/*/*.proto
	@protoc -I . -I /home/_/gh/googleapis/googleapis -I /home/_/gh/bufbuild/protoc-gen-validate \
		--doc_out . \
		--doc_opt markdown,apidocs.md,source_relative \
		proto/*/*/*/*.proto
	@go mod tidy
	@go clean -testcache -fuzzcache -cache

.PHONY: generate_apidocs
generate_apidocs:
	@protoc -I . \
		--doc_out . \
		--doc_opt markdown,apidocs.md,source_relative \
		proto/*/*.proto
	@protoc -I . \
		--doc_out . \
		--doc_opt html,apidocs.html,source_relative \
		proto/*/*.proto
	@protoc -I . \
		--doc_out . \
		--doc_opt json,apidocs.json,source_relative \
		proto/*/*/*/*.proto
	@go mod tidy

.PHONY: go_install_dependencies
go_install_dependencies:
	@go mod tidy
	@go clean -testcache -fuzzcache -cache
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@go install github.com/envoyproxy/protoc-gen-validate@latest
	@go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest

.PHONY: go_install_other
go_install_other:
	@go mod tidy
	@go clean -testcache -fuzzcache -cache
	@go install github.com/99designs/gqlgen@latest
	@go install github.com/martinxsliu/protoc-gen-graphql@latest
	@go install github.com/bufbuild/buf/cmd/buf@latest \
								github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@latest \
								github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@latest

# Firebase Local Emulator Suite
# https://firebase.google.com/docs/emulator-suite
# https://firebase.google.com/docs/emulator-suite/install_and_configure
# https://firebase.google.com/docs/emulator-suite/connect_auth
# 1 install firebase
# 2 init firebase emulators - firebase init emulators 
# 3 start firebase emulators - firebase emulators:start
# 4 inform firebase go sdk to use emulators - export FIREBASE_AUTH_EMULATOR_HOST="localhost:9099"

.PHONY: firebase_init_emulators
firebase_init_emulators:
	@echo "firebase init emulators"
	@firebase init emulators

.PHONY: firebase_emulators_start
firebase_emulators_start:
	@echo "firebase emulators:start"
	@firebase emulators:start
