BINARY_NAME := addr_sh
GO := go
GOFLAGS := -mod=vendor

TOOLS_DIR := $(shell $(GO) env GOPATH)/bin

STATICCHECK := $(TOOLS_DIR)/staticcheck
GOVULNCHECK := $(TOOLS_DIR)/govulncheck
GOSEC := $(TOOLS_DIR)/gosec
OSV_SCANNER := $(TOOLS_DIR)/osv-scanner

.PHONY: all build build-openbsd test check clean tools staticcheck govulncheck gosec osv-scanner help

.DEFAULT_GOAL := help

help:
	@echo "Available targets:"
	@echo "  all           - run check, test, and build"
	@echo "  build         - build binary for current platform"
	@echo "  build-openbsd - build binary for OpenBSD amd64"
	@echo "  test          - run tests"
	@echo "  check         - run all static analysis tools"
	@echo "  staticcheck   - run staticcheck"
	@echo "  govulncheck   - run govulncheck"
	@echo "  gosec         - run gosec"
	@echo "  osv-scanner   - run osv-scanner"
	@echo "  tools         - install all analysis tools"
	@echo "  clean         - remove built binaries"

all: check test build

build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .

build-openbsd:
	GOOS=openbsd GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME)-openbsd-amd64 .

test:
	$(GO) test $(GOFLAGS) ./...

check: staticcheck govulncheck gosec osv-scanner

staticcheck: $(STATICCHECK)
	$(STATICCHECK) ./...

govulncheck: $(GOVULNCHECK)
	$(GOVULNCHECK) ./...

gosec: $(GOSEC)
	$(GOSEC) ./...

osv-scanner: $(OSV_SCANNER)
	$(OSV_SCANNER) scan --lockfile=go.mod

tools: $(STATICCHECK) $(GOVULNCHECK) $(GOSEC) $(OSV_SCANNER)

$(STATICCHECK):
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest

$(GOVULNCHECK):
	$(GO) install golang.org/x/vuln/cmd/govulncheck@latest

$(GOSEC):
	$(GO) install github.com/securego/gosec/v2/cmd/gosec@latest

$(OSV_SCANNER):
	$(GO) install github.com/google/osv-scanner/cmd/osv-scanner@latest

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-openbsd-amd64
