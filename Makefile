PLATFORM                 := $(shell uname)
NAME                     := $(shell basename $(CURDIR))
GOFMT_FILES              ?= $$(find ./ -name '*.go' | grep -v vendor | grep -v externalmodels)
GOTEST_DIRECTORIES       ?= $$(find ./ -type f -iname "*_test.go" -exec dirname {} \; | uniq)

.PHONY: test
test:
	@echo "==> Executing tests..."
	@echo ${GOTEST_DIRECTORIES} | xargs -n1 go test --timeout 30m -v -count 1

.PHONY: swag
swag:
	@echo "==> Updating documentation using swag..."
	@swag init -d cmd/,internal/handler/ --pd --output api/docs/