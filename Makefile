#! /usr/bin/make
#(C) Copyright 2022 Hewlett Packard Enterprise Development LP
# Inspiration from https://github.com/rightscale/go-boilerplate/blob/master/Makefile

NAME=$(shell find cmd -name ".gitkeep_provider" -exec dirname {} \; | sort -u | sed -e 's|cmd/||')

# version shouldn't have 'v' prefix for >= 0.13
VERSION=$(shell cat ./version)
# Change DUMMY_PROVIDER below to reflect the name of the service under development.  The
# value of this variable is used in LOCAL_LOCATION, and is also used in the
DUMMY_PROVIDER=metal
LOCAL_LOCATION=~/.local/share/terraform/plugins/terraform.example.com/$(DUMMY_PROVIDER)/hpegl/$(VERSION)/linux_amd64/

# Stuff that needs to be installed globally (not in vendor)
DEPEND=

# Will get the branch name
SYMBOLIC_REF=$(shell if [ -n "$$CIRCLE_TAG" ] ; then echo $$CIRCLE_TAG; else git symbolic-ref HEAD | cut -d"/" -f 3; fi)
COMMIT_ID=$(shell git rev-parse --verify HEAD)
DATE=$(shell date +"%F %T")

PACKAGE := $(shell git remote get-url origin | sed -e 's|http://||' -e 's|^.*@||' -e 's|.git||' -e 's|:|/|')
VERSION_PACKAGE=$(PACKAGE)/pkg/cmd/$@
VFLAG=-X '$(VERSION_PACKAGE).name=$@' \
      -X '$(VERSION_PACKAGE).version=$(SYMBOLIC_REF)' \
      -X '$(VERSION_PACKAGE).buildDate=$(DATE)' \
      -X '$(VERSION_PACKAGE).buildSha=$(COMMIT_ID)'
TAGS=

# kelog issue: https://github.com/rjeczalik/notify/issues/108
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	TAGS=-tags kqueue
endif
TMPFILE := $(shell mktemp)

GOOSALT ?= 'linux'
GOOS='$(GOOSALT)'

LOCALIZATION_FILES := $(shell find . -name \*.toml | grep -v vendor | grep -v ./bin)

$(NAME): $(shell find . -name \*.go)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=amd64 go build $(TAGS) -ldflags "$(VFLAG)" -o build/$@ .
	
default: all
.PHONY: default

vendor: go.mod go.sum
	go mod download
	go mod vendor

update up: really-clean vendor
.PHONY: update up

clean:
	rm -rf gathered_logs build .vendor/pkg $(testreport_dir) $(coverage_dir)
.PHONY: clean

really-clean clean-all cleanall: clean
	rm -rf vendor
.PHONY: really-clean clean-all cleanall

procs := $(shell grep -c ^processor /proc/cpuinfo 2>/dev/null || echo 1)

go-lint: vendor golangci-lint-config.yaml
	@golangci-lint --version
	golangci-lint run --config golangci-lint-config.yaml --new-from-rev origin/master --max-issues-per-linter 0 --max-same-issues 0
.PHONY: go-lint

tf-lint:
	@terraform fmt -check -recursive -diff || \
		(echo "Some terraform files need to be formatted. Run 'terraform fmt -recursive' to fix them." && return false)
.PHONY: tf-lint

lint: go-lint tf-lint
.PHONY: lint

tools: vendor
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint
.PHONY: tools

testreport_dir := test-reports
test:
	go test -v ./...
.PHONY: test

coverage_dir := coverage/go
coverage: vendor
	@mkdir -p $(coverage_dir)/html
	go test -coverpkg=./... -coverprofile=$(coverage_dir)/coverage.out -v $$(go list ./... | grep -v /vendor/)
	@go tool cover -html=$(coverage_dir)/coverage.out -o $(coverage_dir)/html/main.html;
	@echo "Generated $(coverage_dir)/html/main.html";
.PHONY: coverage

acceptance:	
	TF_ACC=true go test -v -timeout=180s ./...

build: vendor $(NAME)
.PHONY: build

install: build $(NAME)
	# terraform >= v0.13
	mkdir -p $(LOCAL_LOCATION)
	cp build/$(NAME) $(LOCAL_LOCATION)
.PHONY: install

docs: vendor
	go generate ./main.go
.PHONY: docs

all: lint test
.PHONY: all
