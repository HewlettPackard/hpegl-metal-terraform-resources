TEST?=$$(go list ./quake | grep -v vendor)
TF_VERSION:=$(shell cat ./version)

default: build

build: 
	(cd cmd; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o terraform-provider-quake_$(TF_VERSION))

test:
	go test -v -i $(TEST) 
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTGARHS) -timeout=5s

acceptance:
	go test -v -i $(TEST) 
	echo $(TEST) | \
		TF_ACC=true xargs -t -n4 go test -v $(TESTGARHS) -timeout=60s -cover

install: build	
	mkdir -p ~/.terraform/plugins
	cp cmd/terraform-provider-quake* ~/.terraform/plugins
