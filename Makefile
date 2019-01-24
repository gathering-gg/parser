# The version identifier
#
# Default: latest git tag
VERSION ?= $(shell git describe --abbrev=0 --tags)

# The directory to use for build artifacts
#
# Default: "out"
DIST_DIR ?= out

# The binary name without any prefixes
BINARY_NAME := gathering

# The Root URL to make HTTP requests to
#
# Defaults to production
ROOT_URL ?= https://api.gathering.gg

# Go's path to the source code (i.e., in the $GOPATH/src directory)
#
# Default: the current directory within github.com/virtyx-technologies
REPO_PATH ?= github.com/gathering-gg/$(shell basename "$(PWD)")

# Additional flags to pass to the `go build` command. For example, some
# `-ldflags`.
BUILD_FLAGS = -ldflags '-X $(REPO_PATH)/config.Version=$(VERSION) -X $(REPO_PATH)/config.Root=$(ROOT_URL)'


# The go package to build
#
# Default: ./cli
GO_BUILD_TARGET ?= ./cli

# The command that is run to build the program
#
# Default: `go get` followed by `go build`
BUILD_CMD ?= /bin/sh -c "go get -t ./... && go build $(BUILD_FLAGS) -o $@ $(GO_BUILD_TARGET)"

# In case GOOS isn't set, try to figure it out
ifndef GOOS
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	GOOS := darwin
else ifeq ($(UNAME_S),Linux)
	GOOS := linux
else
	GOOS := windows
endif
endif

ifeq ($(OS),Windows_NT) 
    BINARY := $(BINARY_NAME)_$(VERSION)_$(GOOS).exe
else
    BINARY := $(BINARY_NAME)_$(VERSION)_$(GOOS)
endif

$(DIST_DIR):
	mkdir -p $@

release: $(DIST_DIR) $(DIST_DIR)/$(BINARY) 
	cd "$(DIST_DIR)" && zip "$(BINARY).zip" "$(BINARY)"
	rm "$(DIST_DIR)/$(BINARY)"

$(DIST_DIR)/$(BINARY): 
	GOOS=$(GOOS) GOARCH=amd64 $(BUILD_CMD)

build:
	go build -ldflags \
	  "-X 'github.com/gathering-gg/parser/config.Root=http://localhost:8600' -X 'gitlab.com/gathering-gg/gathering.version=$(VERSION)'" \
	  -o gathering \
	  ./cli

test:
	go test -v ./...

cov:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html

clean:
	rm -rf out
	
.PHONY: test clean cov build release
