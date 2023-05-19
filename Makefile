all: run


.PHONY: build
build:
	DOCKER_BUILDKIT=1 docker build -t backend .

.PHONY: run
run: build
	docker compose -f docker-compose.yml up -d --build backend

.PHONY: unit-test
unit-test:
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
	cp docker-compose.yml docker-compose.yml.old
	cp prod/docker-compose.prod.yml docker-compose.yml
	docker compose -f docker-compose.yml up -d --build backend
	cp docker-compose.yml.old docker-compose.yml
	rm docker-compose.yml.old

.PHONY: metrics
metrics:
	cp docker-compose.yml docker-compose.yml.old
	cp prod/docker-compose.prod.yml docker-compose.yml
	docker compose -f docker-compose.yml up -d node_exporter prometheus grafana
	cp docker-compose.yml.old docker-compose.yml
	rm docker-compose.yml.old

.PHONY: build-images
build-images:
	./scripts/build_images.sh

.PHONY: push-images
push-images: build-images
	./scripts/push_images.sh
	
.PHONY: pull-images
pull-images:
	./scripts/pull_images.sh
