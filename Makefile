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

# Go's path to the source code (i.e., in the $GOPATH/src directory)
#
# Default: the current directory within github.com/virtyx-technologies
REPO_PATH ?= github.com/gathering-gg/$(shell basename "$(PWD)")

# Additional flags to pass to the `go build` command. For example, some
# `-ldflags`.
BUILD_FLAGS = -ldflags "-X $(REPO_PATH)/config.Version=$(VERSION) -X $(REPO_PATH)/config.Root=$(ROOT_URL)"


# The go package to build
#
# Default: .
GO_BUILD_TARGET ?= .

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
$(error "$$GOOS is not defined.")
endif
endif


BINARY := $(BINARY_NAME)_$(VERSION)_$(GOOS)

$(DIST_DIR):
	mkdir -p $@

ifeq ($(GOOS),windows)
release: $(DIST_DIR)/$(BINARY) 
	$(call zip_binary,$(BINARY32),.exe)
	$(call zip_binary,$(BINARY64),.exe)
else ifeq ($(GOOS),linux)
release: $(DIST_DIR)/$(BINARY32) $(DIST_DIR)/$(BINARY64) $(DIST_DIR)/$(BINARYARM5) $(DIST_DIR)/$(BINARYARM6) $(DIST_DIR)/$(BINARYARM7) $(DIST_DIR)/$(BINARYARM8) $(ZIP_DIR)
	$(call zip_binary,$(BINARY32),)
	$(call zip_binary,$(BINARY64),)
	$(call zip_binary,$(BINARYARM5),)
	$(call zip_binary,$(BINARYARM6),)
	$(call zip_binary,$(BINARYARM7),)
	$(call zip_binary,$(BINARYARM8),)
	cd $(DIST_DIR) && rm -f $(BINARY_NAME)
else ifeq ($(GOOS),darwin)
release: $(DIST_DIR)/$(BINARY64) $(ZIP_DIR)
	$(call zip_binary,$(BINARY64),)
	cd $(DIST_DIR) && rm -f $(BINARY_NAME)
else
release: $(DIST_DIR)/$(BINARY32) $(DIST_DIR)/$(BINARY64) $(ZIP_DIR)
	$(call zip_binary,$(BINARY32))
	$(call zip_binary,$(BINARY64))
	cd $(DIST_DIR) && rm -f $(BINARY_NAME)
endif


$(DIST_DIR)/$(BINARY): 
	GOOS=$(GOOS) GOARCH=amd64 $(BUILD_CMD)

build:
	go build -ldflags \
	  "-X 'github.com/gathering-gg/parser/config.Root=http://localhost:8600' -X 'gitlab.com/gathering-gg/gathering.version=0.0.3'" \
	  -o gathering \
	  ./cli

test:
	go test -v ./...

cov:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html

release: test cov build
	go build -ldflags \

	  "-X 'github.com/gathering-gg/parser/config.Root=http://localhost:8600' -X 'gitlab.com/gathering-gg/gathering.version=0.0.3'" \
	  -o gathering \
	  ./cli
	
.PHONY: test clean cov build release
