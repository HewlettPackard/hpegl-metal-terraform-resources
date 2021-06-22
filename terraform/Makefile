#(C) Copyright 2020-2021 Hewlett Packard Enterprise Development LP

TEST?=$$(go list ./quake | grep -v vendor)
TF_VERSION:=$(shell cat ./version)

# This is for terraform versions > 0.13
TF_VERSION_NUMBER:=$(shell cat ./version | sed 's/v//')
LOCAL_LOCATION=~/.local/share/terraform/plugins/terraform.example.com/quake/quake/$(TF_VERSION_NUMBER)/linux_amd64/

GOOSALT ?= 'linux'
GOOS='$(GOOSALT)'

default: build

vendor:
	go mod download
	go mod vendor

build: vendor
	(cd cmd; CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=amd64 go build -mod=vendor -o terraform-provider-quake_$(TF_VERSION))

test: vendor
	go test -v -i $(TEST) 
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTGARHS) -timeout=5s

acceptance: vendor
	go test -v -i $(TEST) 
	echo $(TEST) | \
		TF_ACC=true xargs -t -n4 go test -v $(TESTGARHS) -timeout=120s -cover

install: build	
	mkdir -p ~/.terraform.d/plugins
	cp cmd/terraform-provider-quake* ~/.terraform.d/plugins
	mkdir -p $(LOCAL_LOCATION)
	cp cmd/terraform-provider-quake* $(LOCAL_LOCATION)/terraform-provider-quake
