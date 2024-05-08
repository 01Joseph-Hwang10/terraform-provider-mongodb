default: testacc 

OS := $(shell uname)
TARGET := ""

decide-target:
ifeq ($(OS),Linux)
    $(eval TARGET := linux_arm)
else ifeq ($(OS),Darwin)
    $(eval TARGET := darwin_amd64)
else
    $(eval TARGET = windows_amd64)
endif


REGISTRY_HOST := registry.terraform.io
NAMESPACE := 01Joseph-Hwang10
PROVIDER_NAME := terraform-provider-mongodb
VERSION := 0.0.0-unpublished-test

# Build the provider
build:
	$(info ************  Building the Provider Binary  ************)
	go install .

generate:
	$(info ************  Generating Provider Documentation  ************)
	go generate ./...

.PHONY: testacc, clean, show-artifact

TESTPATHS ?= ./...

# Run acceptance tests
#
# There's a problem resolving the provider in the acceptance tests.
# `TF_ACC_PROVIDER_NAMESPACE` is a workaround to make the provider available.
# See the link below for more information and track the issue for updates:
#     https://github.com/hashicorp/terraform-plugin-sdk/issues/1171
testacc: decide-target build
	$(info ************  Running Acceptance Tests  ************)
	DEBUG="$(DEBUG)" \
		EXEC_ROOT="$(shell pwd)" \
		TF_ACC=1 \
		TF_CLI_CONFIG_FILE="$(shell pwd)/terraform.tfrc" \
		TF_ACC_PROVIDER_NAMESPACE="$(NAMESPACE)" \
		go test $(TESTPATHS) \
		-v \
		-timeout 120m \
		-parallel 4 \
		$(TESTARGS) 

ARTIFACT_PATH := $(shell go env GOPATH)/bin/$(PROVIDER_NAME)

show-artifact:
	$(info ************  Showing Artifact Path  ************)
	ls "$(ARTIFACT_PATH)"

clean:
	$(info ************  Cleaning Up Artifacts  ************)
	rm -rf "$(ARTIFACT_PATH)"