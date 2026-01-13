.PHONY: build install version

# Get git info for version
GIT_TAG := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS := -X 'main.version=$(GIT_TAG)' \
           -X 'main.commit=$(GIT_COMMIT)' \
           -X 'main.buildDate=$(BUILD_DATE)'

build:
	@echo "Building git-identity-switcher..."
	@go build -ldflags "$(LDFLAGS)" -o git-identity-switcher .

install:
	@echo "Installing git-identity-switcher..."
	@go install -ldflags "$(LDFLAGS)" .

version:
	@echo "Version: $(GIT_TAG)"
	@echo "Commit: $(GIT_COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

