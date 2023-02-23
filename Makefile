all: run

.PHONY: build
build:
	DOCKER_BUILDKIT=1 docker build -t backend .

.PHONY: run
run: build
	docker compose -f docker-compose.yml up -d --build backend

.PHONY: unit-test
unit-test: db
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
