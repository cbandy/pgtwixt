benchmark:
	go test -bench . -benchmem

.cover.profile: *.go
	go test -coverprofile $@

coverage: .cover.profile
	go tool cover -func $<

coverage-report: .cover.profile
	@command -v gocov > /dev/null || go get github.com/axw/gocov/gocov
	gocov convert $< | gocov annotate - | less -S

setup:
	git submodule update --init --recursive

.PHONY: benchmark coverage coverage-report setup
