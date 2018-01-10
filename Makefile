
behave := .env/bin/behave
pip := .env/bin/pip

.PHONY: benchmark
benchmark: *.go
	go test -bench . -benchmem

.PHONY: check
check: *.go
	go test
	PG_BIN=/usr/lib/postgresql/9.6/bin $(behave)

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
