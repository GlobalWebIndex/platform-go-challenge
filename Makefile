include .env.test
export

.PHONY: build run test clean docker docker-run api-test-all use-user1 use-user2 use-user3 swagger exec-sql

build:
	go build -o bin/challenge ./cmd/challenge

run: build
	./bin/challenge

DOCKER_IMAGE := asset-service:latest
DOCKER_CONTAINER_NAME := asset-service
SQLITE_DB := app.db

docker-run: docker-build
	docker rm -f $(DOCKER_CONTAINER_NAME) || true
	docker run --name $(DOCKER_CONTAINER_NAME) --env-file=.env.test -p 8080:8080 -d $(DOCKER_IMAGE)

docker-build:
	docker build -t $(DOCKER_IMAGE) .

exec-sql:
	docker exec -it $(DOCKER_CONTAINER_NAME) sh -c "sqlite3 $(SQLITE_DB)"

swagger:
	cd internal/api/handlers && swag init \
		--generalInfo handlers.go \
		--output ../docs \
		--parseVendor \
		--parseDependency

define generate_token
	@echo "JWT_SECRET=your-test-secret-key" > .env.test
	@go run tools/jwt/gentoken.go $(1) >> .env.test
	@echo "$(1) set"
endef

use-user1:
	$(call generate_token,user1)

use-user2:
	$(call generate_token,user2)

use-user3:
	$(call generate_token,user3)


api-test-all: api-test-auth api-list api-add api-remove api-update

HOST ?= localhost:8080
PAGE ?= 1
LIMIT ?=20
DESC ?="Updated description"

api-test-auth:
	curl -X GET "http://$(HOST)/" \
		-H "Authorization: Bearer ${JWT_TOKEN}" \
		-H "Content-Type: application/json"

api-list:
	curl -X GET "http://$(HOST)/favorites?page=$(PAGE)&limit=$(LIMIT)" \
		-H "Authorization: Bearer ${JWT_TOKEN}" \
		-H "Content-Type: application/json"

api-add:
	curl -X POST "http://$(HOST)/favorites/$(ASSET_ID)" \
		-H "Authorization: Bearer ${JWT_TOKEN}" \
		-H "Content-Type: application/json"

api-remove:
	curl -X DELETE "http://$(HOST)/favorites/$(ASSET_ID)" \
		-H "Authorization: Bearer ${JWT_TOKEN}" \
		-H "Content-Type: application/json"

api-update:
	curl -X PATCH "http://$(HOST)/favorites/$(ASSET_ID)/description" \
		-H "Authorization: Bearer ${JWT_TOKEN}" \
		-H "Content-Type: application/json" \
		-d '{"description":$(DESC)}'
