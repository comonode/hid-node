#!/usr/bin/make -f

VERSION := $(shell git describe --tags --abbrev=0)
COMMIT := $(shell git rev-parse --short HEAD)

BUILD_DIR ?= $(CURDIR)/build
HIDNODE_CMD_DIR := $(CURDIR)/cmd/hid-noded

GOBIN = $(shell go env GOPATH)/bin
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

SDK_VERSION := $(shell go list -m github.com/cosmos/cosmos-sdk | sed 's:.* ::')
TM_VERSION := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=hid-node \
	-X github.com/cosmos/cosmos-sdk/version.AppName=hid-node \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X github.com/tendermint/tendermint/version.TMCoreSemVer=$(TM_VERSION)

BUILD_FLAGS := -ldflags '$(ldflags)'

export GO111MODULE=on

###############################################################################
###                                  Build                                  ###
###############################################################################
.PHONY: build install 

all: proto-gen swagger-docs-gen build

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		@go mod verify

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) $(HIDNODE_CMD_DIR)	

build:
	go build -mod=readonly $(BUILD_FLAGS) -o $(BUILD_DIR)/hid-noded $(HIDNODE_CMD_DIR)

###############################################################################
###                                  Proto                                  ###
###############################################################################

proto-gen:
	@echo "Generating golang code from protobuf"
	./scripts/protocgen.sh

###############################################################################
###                                  Docs                                   ###
###############################################################################

swagger-docs-gen:
	@echo "Generating swagger docs"
	./scripts/protoc-swagger-gen.sh
