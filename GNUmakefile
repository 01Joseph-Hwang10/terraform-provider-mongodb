default: testacc

# Build the provider
build:
	go build -o dist/terraform-provider-mongodb

TARGET ?= ./...

# Run acceptance tests
.PHONY: testacc
testacc: build
	TF_CLI_CONFIG_FILE="$(shell pwd)/.terraform.tfrc" TF_ACC=1 go test $(TARGET) -v $(TESTARGS) -timeout 120m
