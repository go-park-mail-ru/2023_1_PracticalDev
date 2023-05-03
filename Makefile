all: run


.PHONY: build
build:
	DOCKER_BUILDKIT=1 docker build -t backend .

.PHONY: run
run: build
	docker compose -f docker-compose.yml up -d --build backend

.PHONY: unit-test
unit-test: db redis
	DOCKER_BUILDKIT=1 docker build . --target unit-test --network=host

.PHONY: unit-test-with-coverage
unit-test-with-coverage: unit-test
	DOCKER_BUILDKIT=1 docker build . --target coverage --output ./ --network=host

.PHONY: print-coverage
print-coverage: unit-test
	@DOCKER_BUILDKIT=1 docker build . --target print-coverage --network=host -t coverage
	docker run --rm coverage
	@docker image rm coverage >/dev/null

.PHONY: lint
lint:
	DOCKER_BUILDKIT=1 docker build . --target lint

.PHONY: db
db:
	docker compose -f docker-compose.yml up -d db

.PHONY: redis
redis:
	docker compose -f docker-compose.yml up -d redis

.PHONY: migrations
migrations:
	./scripts/run_migrations.sh

.PHONY: fill-test-data
fill-test-data:
	./scripts/run_migrations.sh
	./scripts/populate_db.sh

.PHONY: mocks
mocks:
	./scripts/gen_mocks.sh

.PHONY: metrics-test
metrics-test:
	docker compose -f docker-compose.yml up -d node_exporter prometheus grafana

.PHONY: deploy
deploy:
	docker compose -f prod/docker-compose.prod.yml up -d --build backend

.PHONY: metrics
metrics:
	docker compose -f prod/docker-compose.prod.yml up -d node_exporter prometheus grafana