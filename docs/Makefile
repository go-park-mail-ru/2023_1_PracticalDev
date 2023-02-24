all: lint

.PHONY: lint
lint:
	@docker run --rm -t -v "$(PWD)":/tmp stoplight/spectral lint -v --ruleset "/tmp/.spectral.yml" "/tmp/openapi.yaml"
