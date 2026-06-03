.PHONY: help test test-race test-cover clean

GO ?= go
PKGS ?= ./...
TESTFLAGS ?=
COVERPROFILE ?= coverage.out
GOCACHE ?= $(CURDIR)/.cache/go-build

help:
	@printf "Available targets:\n"
	@printf "  make test        Run all tests\n"
	@printf "  make test-race   Run all tests with race detector\n"
	@printf "  make test-cover  Run all tests with coverage profile\n"
	@printf "  make clean       Remove generated test artifacts\n"

test:
	GOCACHE=$(GOCACHE) $(GO) test $(TESTFLAGS) $(PKGS)

test-race:
	GOCACHE=$(GOCACHE) $(GO) test -race $(TESTFLAGS) $(PKGS)

test-cover:
	GOCACHE=$(GOCACHE) $(GO) test -coverprofile=$(COVERPROFILE) $(TESTFLAGS) $(PKGS)

clean:
	rm -f $(COVERPROFILE)
	rm -rf $(GOCACHE)
