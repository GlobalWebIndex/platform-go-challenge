.DEFAULT_GOAL := help
DOCKER_COMPOSE_YAML := infra/deploy/local/docker-compose.yml

# `make help` generates a help message for each target that
# has a comment starting with ##
help:
	@echo "Please use 'make <target>' where <target> is one of the following:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

serve: ## Serve the app with Docker Compose
	docker-compose -f $(DOCKER_COMPOSE_YAML) up -d --build --force-recreate  platform-go-challenge

stop: ## Stop the Docker Compose app
	docker-compose -f $(DOCKER_COMPOSE_YAML) down

lint: ## Perform linting
	golint -set_exit_status=1 `go list ./...`

test: ## Run unit tests
	go test ./... -race -cover -tags=integration

ci: ## Run the CI pipeline
	docker-compose -f $(DOCKER_COMPOSE_YAML) build platform-go-challenge_ci
	docker-compose -f $(DOCKER_COMPOSE_YAML) run platform-go-challenge_ci ./script/ci.sh
