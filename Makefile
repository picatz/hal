# The name of the binary to build
NAME := hal

# The version of the binary to build, from the git tag
VERSION := $(shell git describe --tags --always --dirty)

# The git commit hash to embed in the binary
GIT_COMMIT := $(shell git rev-parse HEAD)

# The date to embed in the binary
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# The go build command
GO_BUILD := go build -ldflags "-X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT) -X main.date=$(DATE)"

# The go install command
# GO_INSTALL := go install -ldflags "-X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT) -X main.date=$(DATE)"

# The build target
.PHONY: build
build:
	$(GO_BUILD) -o $(NAME)

# Record the demo GIF
.PHONY: demo
demo:
	vhs < demo.tape