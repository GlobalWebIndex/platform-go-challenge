start:
	make stop
	docker compose -f ./deploy/docker-compose.yaml build
	docker compose -f ./deploy/docker-compose.yaml up --detach

stop:
	docker compose -f ./deploy/docker-compose.yaml down

migration-up:
	make _internal_run-migration ACTION=up
migration-down:
	make _internal_run-migration ACTION=down

_internal_run-migration:
	@docker build \
		--tag migrations \
		-f deploy/development/migrations.dockerfile .
	@docker run \
		--network=gwi \
		--env ACTION=$(ACTION) \
		migrations
