TEST?=$$(go list ./quake | grep -v vendor)

default: build

build: 
	(cd cmd; ./build.sh)

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
