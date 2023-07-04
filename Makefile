appname ?= app-gwi
apptag = localhost/$(appname):latest

source_dir := $(abspath $(dir $(lastword ${MAKEFILE_LIST})))

# for container use podman or docker
container ?= podman
# Containerfile Dockerfile - are the same for now
containerfile ?= Containerfile

aqlV = 3.11.1
dbaql = app-dbaql-$(aqlV)
# 3.11.1
# run -e ARANGO_NO_AUTH=1 

cqlV = 5.3.0-rc0
dbcql = app-dbcql-$(cqlV)
# 5.2.3 5.3.0-rc0

sqlV = 15.3
dbsql = app-dbsql-$(sqlV)
# 15.3 16beta1
# run -e POSTGRES_PASSWORD=trust 
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= postgres

ip = $(shell hostname -i)
app_host = $(shell hostname -i)
HOSTNAME ?= $(shell hostname)
APP_NAME ?= $(appname)
APP_MODE ?= ""
SERVER_GRPC_ADDRESS ?= ":9090"
SERVER_RESTGW_ADDRESS ?= ":9080"
APP_STORAGE_AQL_USERNAME ?= "arango"
APP_STORAGE_AQL_PASSWORD ?= "arango"
APP_STORAGE_AQL_ENDPOINTS ?= ""

.PHONY: go-all
go-all: go-clean go-tidy go-fmt go-vuln go-lint-clean-cache go-lint

.PHONY: go-tidy
go-tidy:
	@go mod tidy

.PHONY: go-fmt
go-fmt:
	@go fmt ./...

.PHONY: go-vuln
go-vuln:
	@go mod tidy
	@govulncheck ./...

.PHONY: go-lint-clean-cache
go-lint-clean-cache:
	@golangci-lint cache clean

.PHONY: go-lint
go-lint:
	@go mod tidy
	@golangci-lint run

.PHONY: go-lint-v
go-lint-v:
	@go mod tidy
	@golangci-lint run -v > /dev/null

.PHONY: go-lint-x
go-lint-x:
	@go mod tidy
	@golangci-lint run --no-config --enable-all -D=gofumpt,gci,varnamelen,exhaustivestruct,structcheck,nosnakecase,golint,maligned,scopelint,ifshort,varcheck,interfacer,deadcode,depguard

.PHONY: go-clean
go-clean:
	@go mod tidy
	@go clean -testcache -fuzzcache -cache

.PHONY: go-test-example-x
go-test-example-x:
	@go mod tidy
	@go test -timeout 90s -run ^ExampleX$

.PHONY: go-test-x
go-test-x:
	@go mod tidy
	@go clean -testcache -fuzzcache
	@go test -timeout 90s -run ^TestX$

.PHONY: go-test
go-test:
	@go mod tidy
	@go clean -testcache -fuzzcache
	@go test ./test

.PHONY: go-test-race
go-test-race:
	@go mod tidy
	@go clean -testcache -fuzzcache
	@go test -race ./test

.PHONY: go-bench-x
go-bench-x:
	@go mod tidy
	@go clean -testcache -fuzzcache
	@go test -benchmem -bench .

.PHONY: go-run-main-server
go-run-main-server:
	@go mod tidy
	@-go run .

.PHONY: generate-protoc
generate-protoc:
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

.PHONY: generate-apidocs
generate-apidocs:
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

.PHONY: go-install-dependencies
go-install-dependencies:
	@go mod tidy
	@go clean -testcache -fuzzcache -cache
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@go install github.com/envoyproxy/protoc-gen-validate@latest
	@go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest

.PHONY: go-install-other
go-install-other:
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

.PHONY: firebase-init-emulators
firebase-init-emulators:
	@echo "firebase init emulators"
	@firebase init emulators

.PHONY: firebase-emulators-start
firebase-emulators-start:
	@echo "firebase emulators:start"
	@firebase emulators:start

# --env APP_HOST=${APP_HOST}
# set local host ip into url with port http://localhost:8529,
# APP_STORAGE_AQL_ENDPOINTS ?= "http://192.168.1.15:8529,http://localhost:8529,http://app-dbaql-3.11.1:8529"
# @echo pods: "http://$(dbaql):8529,http://localhost:8529,http://$(dbaql)-podified:8529"
# sudo netstat -tulpn | grep :9080

.PHONY: container-0-check-host
container-0-check-host:
	@echo APP_NAME:		$(appname)
	@echo container:	$(container)
	@echo HOSTNAME:		$(HOSTNAME)
	@echo host_ips:		$(app_host)
	@echo dbaql:		$(dbaql)
	@echo dbcql:		$(dbcql)
	@echo dbsql:		$(dbsql)
	@$(container) container ps

.PHONY: container-0-list-containers
container-0-list-containers:
	@$(container) container ls -a
	@echo
	@$(container) image ls
	@echo
	@$(container) system df
	@echo -e "\nOK\n"

.PHONY: container-1-run-db-aql-once
container-1-run-db-aql-once:
	@echo "preparing $(dbaql)"
	@until $(container) run -d \
	--name $(dbaql) \
	--hostname arango \
	-e ARANGO_NO_AUTH=1 \
	-p 8529:8529 \
	docker.io/arangodb/arangodb:$(aqlV) \
	; do sleep 10; done
	@echo
	@$(container) ps
	@echo -e "\nOK"
	@echo -e "\n$(dbaql) - prepared and started. you can stop it by: make container-2-stop-db-aql"
	@echo -e "open AQL UI http://localhost:8529/\n"

.PHONY: container-2-stop-db-aql
container-2-stop-db-aql:
	@echo "stoping"
	@until $(container) stop $(dbaql); do sleep 10; done
	@$(container) container ls
	@echo -e "\n$(dbaql) - stopped. you can start it by: make container-3-start-db-aql\n"

.PHONY: container-3-start-db-aql
container-3-start-db-aql:
	@echo "starting"
	@until $(container) start $(dbaql); do sleep 10; done
	@$(container) container ls
	@echo -e "\n$(dbaql) - started. you can stop it by: make container-2-stop-db-aql"
	@echo -e "open AQL UI http://localhost:8529/\n"

.PHONY: container-1-run-db-cql-once
container-1-run-db-cql-once:
	@echo "preparing $(dbcql)"
	@until $(container) run -d \
	--name $(dbcql) \
	--hostname scylla \
	-p 10000:10000 \
	-p 19042:19042 \
	-p 19142:19142 \
	-p 9042:9042 \
	-p 9142:9142 \
	-p 7000:7000 \
	-p 7001:7001 \
	-p 9180:9180 \
	-p 9100:9100 \
	docker.io/scylladb/scylla:$(cqlV) \
	--smp 4 \
	--developer-mode 1 \
	--experimental 1 \
	--api-address 0.0.0.0 \
	; do sleep 10; done
	@echo
	@$(container) ps
	@echo -e "\nOK"
	@echo -e "\n$(dbcql) - prepared and started. you can stop it by: make container-2-stop-db-cql"
	@echo -e "open CQL UI http://localhost:10000/\n"

# --listen-address 0.0.0.0 \
# --rpc-address 0.0.0.0 \
# --api-address 0.0.0.0 \
# --broadcast-rpc-address 0.0.0.0 \
# --network=host \

.PHONY: container-2-stop-db-cql
container-2-stop-db-cql:
	@echo "stoping"
	@until $(container) stop $(dbcql); do sleep 10; done
	@$(container) container ls
	@echo -e "\n$(dbcql) - stopped. you can start it by: make container-3-start-db-cql\n"

.PHONY: container-3-start-db-cql
container-3-start-db-cql:
	@echo "starting"
	@until $(container) start $(dbcql); do sleep 10; done
	@$(container) container ls
	@echo -e "\n$(dbcql) - started. you can stop it by: make container-2-stop-db-cql"
	@echo -e "open CQL UI http://localhost:10000/\n"

.PHONY: container-1-run-db-sql-once
container-1-run-db-sql-once:
	@echo "preparing $(dbsql)"
	@until $(container) run -d \
	--name $(dbsql) \
	--hostname postgres \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-p 5432:5432 \
	docker.io/library/postgres:$(sqlV) \
	; do sleep 10; done
	@echo
	@$(container) ps
	@echo -e "\nOK"
	@echo -e "\n$(dbsql) - prepared and started. you can stop it by: make container-2-stop-db-sql"
	@echo -e "open SQL at localhost:5432\n"

# -e POSTGRES_PASSWORD=trust \
# -c listen_addresses='*' \

.PHONY: container-2-stop-db-sql
container-2-stop-db-sql:
	@echo "stoping"
	@until $(container) stop $(dbsql); do sleep 10; done
	@$(container) container ls
	@echo -e "\n$(dbsql) - stopped. you can start it by: make container-3-start-db-sql\n"

.PHONY: container-3-start-db-sql
container-3-start-db-sql:
	@echo "starting"
	@until $(container) start $(dbsql); do sleep 10; done
	@$(container) container ls
	@echo -e "\n$(dbsql) - started. you can stop it by: make container-2-stop-db-sql"
	@echo -e "open SQL UI http://localhost:8529/\n"

# docker build -t localhost/app-gwi:latest .
# docker build -f Containerfile .

.PHONY: container-4-build-app
container-4-build-app:
	@echo "*** container-build - $(appname) ***"
	@$(container) build --no-cache \
	-t $(appname) \
	-f $(containerfile) .
	@echo -e "\nOK\n"
	@$(container) image prune -f
	@$(container) volume prune -f
	@echo
	@$(container) image ls

.PHONY: container-5-run-app--rm
container-5-run-app--rm:
	@$(container) container ps
	@echo -e "\npassing host_ip to container $(app_host)\n" 
	@-$(container) run --rm \
	--name $(appname) \
	--hostname $(appname) \
	-p 9090:9090 \
	-p 9080:9080 \
	-e APP_NAME=$(APP_NAME) \
	-e APP_MODE=$(APP_MODE) \
	-e APP_HOST='$(app_host)' \
	-e SERVER_GRPC_ADDRESS=$(SERVER_GRPC_ADDRESS) \
	-e SERVER_RESTGW_ADDRESS=$(SERVER_RESTGW_ADDRESS) \
	-e APP_STORAGE_AQL_USERNAME=$(APP_STORAGE_AQL_USERNAME) \
	-e APP_STORAGE_AQL_PASSWORD=$(APP_STORAGE_AQL_PASSWORD) \
	-e APP_STORAGE_AQL_ENDPOINTS=$(APP_STORAGE_AQL_ENDPOINTS) \
	$(appname)
	@echo -e "\nOK\n"
	@$(container) image prune -f
	@$(container) container ps
	@$(container) image ls

# --requires is supported by podman but not by docker
# --requires $(dbaql) \

.PHONY: container-x-go-test-Example_gRPC_Client_loading_fake_data
container-x-go-test-Example_gRPC_Client_loading_fake_data:
	@go mod tidy
	@go clean -testcache -fuzzcache -cache
	@go test -v -timeout 5m $(source_dir)/test/fake

.PHONY: container-x-run-app-once-no--rm
container-x-run-app-once-no--rm:
	@$(container) container ps
	@echo "*** container-run ***"
	@-$(container) run \
	--name $(appname) \
	--hostname $(appname) \
	-p 9090:9090 \
	-p 9080:9080 \
	-e APP_NAME=$(APP_NAME) \
	-e APP_MODE=$(APP_MODE) \
	-e APP_HOST='$(app_host)' \
	-e SERVER_GRPC_ADDRESS=$(SERVER_GRPC_ADDRESS) \
	-e SERVER_RESTGW_ADDRESS=$(SERVER_RESTGW_ADDRESS) \
	-e APP_STORAGE_AQL_USERNAME=$(APP_STORAGE_AQL_USERNAME) \
	-e APP_STORAGE_AQL_PASSWORD=$(APP_STORAGE_AQL_PASSWORD) \
	-e APP_STORAGE_AQL_ENDPOINTS=$(APP_STORAGE_AQL_ENDPOINTS) \
	$(appname)
	@echo -e "\nOK\n"
	@$(container) image prune -f
	@$(container) container ps
	@$(container) image ls

.PHONY: container-x-stop-app
container-x-stop-app:
	@echo "stoping"
	@until $(container) stop $(appname); do sleep 10; done
	@$(container) container ls
	@echo -e "\n$(appname) - stopped. you can start it by: make container-7-start-app\n"

.PHONY: container-x-start-app
container-x-start-app:
	@echo "starting"
	@until $(container) start $(appname); do sleep 10; done
	@$(container) container ls
	@echo -e "\n$(appname) - started. you can stop it by: make container-6-stop-app\n"

.PHONY: container-0-prune
container-0-prune:
	@echo "*** container-prune ***"
	@$(container) image prune -f
	@$(container) volume prune -f
	@$(container) container ps -a
	@$(container) image ls

.PHONY: container-x-remove-app
container-x-remove-app:
	@echo "*** container-remove ***"
	@$(container) image ls
	@echo
	@$(container) container rm $(appname)
	@echo -e "\nOK"
	@$(container) image rm $(appname)
	@echo -e "\nOK"
	@$(container) image prune -f
	@$(container) volume prune -f
	@echo
	@$(container) image ls

.PHONY: container-x-remove-db-aql
container-x-remove-db-aql:
	@echo "*** container-remove ***"
	@$(container) image ls
	@echo
	@$(container) container rm $(dbaql)
	@echo -e "\nOK"
	@$(container) image prune -f
	@$(container) volume prune -f
	@echo
	@$(container) image ls

.PHONY: container-x-nodetool-status-db-cql
container-x-nodetool-status-db-cql:
	@echo "*** $(container) exec -it $(dbcql) nodetool status ***"
	@until $(container) exec -it $(dbcql) nodetool status; do sleep 10; done

.PHONY: container-x-net-seeds-db-cql
container-x-net-seeds-db-cql:
	@echo "***  ***"
	@$(container) inspect $(dbcql) -f '{{ .NetworkSettings.IPAddress }}'
	@$(container) inspect $(dbcql) -f '{{ .NetworkSettings.Ports }}'

# @$(container) inspect $(dbcql) -f '{{ .NetworkSettings.Ports }}'
# @$(container) inspect $(dbcql) -f '{{ .NetworkSettings.IPAddress }} {{ .NetworkSettings.Ports }}'

# podman inspect mysql -f '{{ .NetworkSettings.IPAddress }} {{ .NetworkSettings.Ports }}'
# "$(docker inspect --format='{{ .NetworkSettings.IPAddress }}' some-scylla)"

.PHONY: container-x-remove-db-cql
container-x-remove-db-cql:
	@echo "*** container-remove ***"
	@$(container) image ls
	@echo
	@$(container) container rm $(dbcql)
	@echo -e "\nOK"
	@$(container) image prune -f
	@$(container) volume prune -f
	@echo
	@$(container) image ls

.PHONY: container-x-remove-db-sql
container-x-remove-db-sql:
	@echo "*** container-remove ***"
	@$(container) image ls
	@echo
	@$(container) container rm $(dbsql)
	@echo -e "\nOK"
	@$(container) image prune -f
	@$(container) volume prune -f
	@echo
	@$(container) image ls
