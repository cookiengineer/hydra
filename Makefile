# Hydra build system
#
# Binaries are written to ./build/ with distribution, operating-system and
# architecture suffixes derived from Go's platform tags, for example:
#
#   build/hydra-debian-linux-amd64
#   build/hydra-archlinux-linux-amd64
#   build/hydra-native-linux-amd64
#
# GOOS/GOARCH default to the host's Go values and can be overridden, e.g.
# `make debian GOARCH=arm64` once a cross toolchain is available.

SHELL := /bin/bash

GOOS   ?= $(shell go env GOOS 2>/dev/null || echo linux)
GOARCH ?= $(shell go env GOARCH 2>/dev/null || echo amd64)

BUILD_DIR := build

DEBIAN_IMAGE := hydra-debian-build
ARCH_IMAGE   := hydra-archlinux-build

DEBIAN_BINARY := $(BUILD_DIR)/hydra-debian-$(GOOS)-$(GOARCH)
ARCH_BINARY   := $(BUILD_DIR)/hydra-archlinux-$(GOOS)-$(GOARCH)
NATIVE_BINARY := $(BUILD_DIR)/hydra-native-$(GOOS)-$(GOARCH)

SOURCES := $(shell find . -name '*.go' -not -path './build/*' -not -path './toolchain/*')

CONTAINERFLAGS := \
	-v "$(CURDIR)":/src:ro,Z \
	-v "$(CURDIR)/$(BUILD_DIR)":/build:Z \
	-w /src \
	-e GOOS=$(GOOS) -e GOARCH=$(GOARCH) -e CGO_ENABLED=1

HOST ?= hydratwo

.PHONY: all debian archlinux native unit integration e2e clean help

all: debian archlinux

debian: $(DEBIAN_BINARY)

$(DEBIAN_BINARY): $(SOURCES) go.mod debian.Containerfile | $(BUILD_DIR)
	podman build -t $(DEBIAN_IMAGE) -f debian.Containerfile .
	podman run --rm $(CONTAINERFLAGS) $(DEBIAN_IMAGE) \
		go build -o /build/$(notdir $@) ./cmds/hydra

archlinux: $(ARCH_BINARY)

$(ARCH_BINARY): $(SOURCES) go.mod archlinux.Containerfile | $(BUILD_DIR)
	podman build -t $(ARCH_IMAGE) -f archlinux.Containerfile .
	podman run --rm $(CONTAINERFLAGS) $(ARCH_IMAGE) \
		go build -o /build/$(notdir $@) ./cmds/hydra

native: $(NATIVE_BINARY)

$(NATIVE_BINARY): $(SOURCES) go.mod | $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -o $(NATIVE_BINARY) ./cmds/hydra

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

unit:
	go test ./...

integration: | $(BUILD_DIR)
	podman build -t $(DEBIAN_IMAGE) -f debian.Containerfile .
	podman run --rm \
		-v "$(CURDIR)":/src:ro,Z \
		-w /src \
		-e DISPLAY=:99 \
		$(DEBIAN_IMAGE) \
		sh -c 'Xvfb :99 -screen 0 1920x1080x24 & sleep 1 && go test -tags=integration -v -count=1 ./adapters/xorg/...'

e2e:
	go run ./toolchain/test.go $(HOST) server

clean:
	rm -rf $(BUILD_DIR)

help:
	@echo "Targets:"
	@echo "  make debian        build $(DEBIAN_BINARY)"
	@echo "  make archlinux     build $(ARCH_BINARY)"
	@echo "  make native        build $(NATIVE_BINARY)"
	@echo "  make all           build debian + archlinux"
	@echo "  make unit          run unit tests"
	@echo "  make integration   run Xvfb integration tests (debian container)"
	@echo "  make e2e HOST=...  run two-machine end-to-end tests"
	@echo "  make clean         remove $(BUILD_DIR)/"
