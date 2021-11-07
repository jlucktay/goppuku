# Inspiration:
# - https://devhints.io/makefile
# - https://tech.davis-hansson.com/p/make/

SHELL := bash
.DELETE_ON_ERROR:
.ONESHELL:
.SHELLFLAGS := -euo pipefail -c

MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --warn-undefined-variables

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later.)
endif
.RECIPEPREFIX = >

# Default - top level rule is what gets run when you run just `make` without specifying a target.
build: out/image-id
.PHONY: build

# Clean up the built binary and output directories.
# All the sentinel files go under tmp, so this will cause everything to get rebuilt.
clean:
> rm -rf dist goppuku tmp out
.PHONY: clean

# Tests - re-run if any Go files have changes since tmp/.tests-passed.sentinel last touched.
tmp/.tests-passed.sentinel: $(shell find . -type f -iname "*.go")
> mkdir -p $(@D)
> go test -cover -race -v ./...
> touch $@

# Lint - re-run if the tests have been re-run (and so, by proxy, whenever the source files have changed).
tmp/.linted.sentinel: tmp/.tests-passed.sentinel
> mkdir -p $(@D)
> golangci-lint run
> go vet ./...
> touch $@

# Docker image - re-build if the lint output is re-run.
out/image-id: tmp/.linted.sentinel
> mkdir -p $(@D)
> image_id="go.jlucktay.dev/goppuku:$$(uuidgen)"
> docker build --tag="$${image_id}" .
> echo "$${image_id}" > out/image-id
