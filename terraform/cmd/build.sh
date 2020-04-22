#!/bin/bash
# Copyright (c) 2017-2020 Hewlett Packard Enterprise Development LP.
set -e


TF_VERSION=$(cat ../version)
echo "building terraform-provider-quake_${TF_VERSION}"
# Creaet a versioned build for the REST API to work and  make the executable statically linked. When
# embeding the provider in some containers there are missing shared libraries if dynamically linked.
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o terraform-provider-quake_${TF_VERSION}
#upx terraform-provider-quake_${TF_VERSION}
