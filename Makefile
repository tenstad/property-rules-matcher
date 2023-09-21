GOFUMPT_VERSION ?= v0.5.0

all: fmt test

.PHONY: fmt
fmt:
	go run mvdan.cc/gofumpt@$(GOFUMPT_VERSION) -l -w .

.PHONY: test
test:
	go test -v ./...
