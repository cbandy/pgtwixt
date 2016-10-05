benchmark: *.go
	go test -bench . -benchmem

check: *.go
	go test

.cover.profile: *.go
	go test -coverprofile $@

coverage: .cover.profile
	go tool cover -func $<

coverage-report: .cover.profile
	@command -v gocov > /dev/null || go get github.com/axw/gocov/gocov
	gocov convert $< | gocov annotate - | less -S

setup: .gitmodules
	git submodule update --init --recursive

.PHONY: benchmark check coverage coverage-report setup
