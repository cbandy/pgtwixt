.DEFAULT_GOAL := help

DOCKER_RELEASE ?= stretch
DOCKER_GROUP := $(shell id -g)
DOCKER_USER  := $(shell id -u)

PG_BIN     ?= /usr/lib/postgresql/$(PG_VERSION)/bin
PG_VERSION ?= $(shell ls /usr/lib/postgresql | sort --version-sort | tail --lines=1)

export DOCKER_GROUP DOCKER_RELEASE DOCKER_USER

.PHONY: benchmark
benchmark: *.go
	go test -bench . -benchmem

.PHONY: check
check: *.go features/* radish/* requirements.txt
	go test . ./cmd/pgtwixt
	go build -o radish/pgtwixt ./cmd/pgtwixt
	@command -v radish > /dev/null || pip install --requirement requirements.txt
	PG_BIN='$(PG_BIN)' radish --no-line-jump --with-traceback features

.cover.profile: *.go
	go test -coverprofile $@

.PHONY: coverage
coverage: .cover.profile
	go tool cover -func $<

.PHONY: coverage-report
coverage-report: .cover.profile
	@command -v gocov > /dev/null || go get github.com/axw/gocov/gocov
	gocov convert $< | gocov annotate - | less -S

.PHONY: docker-check
docker-check: ## Run all the tests in a new Docker container
	docker-compose run --rm dev make check

.PHONY: docker-clean
docker-clean: ## Remove any Docker images or volumes created by this project
docker-clean: | docker-clean-volumes docker-clean-images

.PHONY: docker-clean-images
docker-clean-images:
	FOUND="$$(docker images --filter 'label=project=pgtwixt' --quiet)"; [ -z "$$FOUND" ] || docker rmi $$FOUND

.PHONY: docker-clean-volumes
docker-clean-volumes:
	docker-compose down --volumes
	FOUND="$$(docker volume ls --filter 'label=project=pgtwixt' --quiet)"; [ -z "$$FOUND" ] || docker volume rm $$FOUND

.PHONY: docker-dev
docker-dev: ## Start a shell with all the tools needed to run tests
	docker-compose build dev
	docker-compose run --rm dev || true

.PHONY: help
help: ALIGN=16
help: ## Print this message
	@awk -F ': ## ' -- "/^[^':]+: ## /"' { printf "'$$(tput bold)'%-$(ALIGN)s'$$(tput sgr0)' %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

requirements.txt: requirements.in
	@command -v pip-compile > /dev/null || pip install pip-tools
	pip-compile --verbose
