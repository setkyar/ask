# Set the version flag
VERSION = 1.0.0

# Set the project name
PROJECT = ask

# Set the Go build flags
GO_BUILD_FLAGS = -ldflags "-s -w"

# Set the output directory for the binaries
BIN_DIR = bin

# Set the default target
.DEFAULT_GOAL = all

# Build all the binaries
all: build-macos build-linux build-windows

# Build the macOS binary
build-macos:
	GOOS=darwin GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(PROJECT)-darwin-amd64-$(VERSION)

# Build the Linux binary
build-linux:
	GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(PROJECT)-linux-amd64-$(VERSION)

# Build the Windows binary
build-windows:
	GOOS=windows GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(PROJECT)-windows-amd64-$(VERSION).exe

build-docker:
	docker build -t setkyar/$(PROJECT):$(VERSION) .

docker-push:
	docker push setkyar/$(PROJECT):$(VERSION)

run-docker:
	docker run -it setkyar/$(PROJECT):$(VERSION) /bin/bash