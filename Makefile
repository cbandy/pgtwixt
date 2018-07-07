
pip := .env/bin/pip
radish := .env/bin/radish

.PHONY: benchmark
benchmark: *.go
	go test -bench . -benchmem

.PHONY: check
check: *.go features/* radish/*
	go test . ./cmd/pgtwixt
	go build -o radish/pgtwixt ./cmd/pgtwixt
	PG_BIN=/usr/lib/postgresql/9.6/bin $(radish) --no-line-jump --with-traceback features

.cover.profile: *.go
	go test -coverprofile $@

.PHONY: coverage
coverage: .cover.profile
	go tool cover -func $<

.PHONY: coverage-report
coverage-report: .cover.profile
	@command -v gocov > /dev/null || go get github.com/axw/gocov/gocov
	gocov convert $< | gocov annotate - | less -S

.env: requirements.txt
	virtualenv --python=python3 .env
	$(pip) install -r requirements.txt
	@touch $@

.PHONY: setup
setup: .env .gitmodules
	git submodule update --init --recursive
