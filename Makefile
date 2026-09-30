BINARY_NAME ?= go-web-starter
VARS_PKG ?= github.com/SisyphusSQ/go-web-starter/v2/vars
BUILD_DIR ?= bin
GO ?= go
CGO_ENABLED ?= 0
PLATFORMS := windows-amd64 windows-arm64 darwin-amd64 darwin-arm64 linux-amd64 linux-arm64
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

BUILD_FLAGS  = -X '$(VARS_PKG).AppName=$(BINARY_NAME)'
BUILD_FLAGS += -X '$(VARS_PKG).AppVersion=$(VERSION)'
BUILD_FLAGS += -X '$(VARS_PKG).GoVersion=$(shell $(GO) env GOVERSION)'
BUILD_FLAGS += -X '$(VARS_PKG).BuildTime=$(shell date +"%Y-%m-%d %H:%M:%S")'
BUILD_FLAGS += -X '$(VARS_PKG).GitCommit=$(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)'
BUILD_FLAGS += -X '$(VARS_PKG).GitRemote=$(shell git config --get remote.origin.url 2>/dev/null || echo unknown)'

LDFLAGS = -ldflags="$(BUILD_FLAGS)"

.PHONY: help all build release run test lint fmt tidy clean install

help:
	@echo "  build-all / release-all - 构建 Windows、macOS、Linux 的 amd64 和 arm64"
	@echo "  build-<os>-<arch> / release-<os>-<arch> - 单个平台，例如 build-windows-arm64"
	@echo "  跨平台输出：$(BUILD_DIR)/<os>-<arch>/$(BINARY_NAME)[.exe]；不执行测试或发布"
	@echo "Available targets:"
	@echo "  all      - Run fmt, test, and build"
	@echo "  build    - Build local binary"
	@echo "  release  - Build release binary with trimpath"
	@echo "  run      - Run starter locally"
	@echo "  test     - Run tests with race detector"
	@echo "  lint     - Run golangci-lint"
	@echo "  fmt      - Run go fmt"
	@echo "  tidy     - Run go mod tidy"
	@echo "  install  - Install binary into GOPATH/bin"
	@echo "  clean    - Remove build artifacts"

all: fmt test build

build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./main.go

release:
	@mkdir -p $(BUILD_DIR)
	$(GO) build -trimpath $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./main.go

run:
	$(GO) run ./main.go

test:
	$(GO) test -race ./...

lint:
	$(GO) run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0 run

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

install:
	$(GO) install $(LDFLAGS) .

clean:
	rm -rf $(BUILD_DIR)

.PHONY: verify integration vuln
verify:
	@test -z "$$(gofmt -l cmd internal/scaf_fold/*.go main.go vars)"
	$(GO) vet ./...
	$(GO) build ./...

integration:
	$(GO) test -tags integration ./internal/scaf_fold -run TestGenerateE2EDBCombosIntegration -count=1 -timeout=20m

vuln:
	$(GO) run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...

# 保留本机构建入口；六平台目标各自输出到独立目录，允许 make -j 并行。
.PHONY: build-all release-all $(addprefix build-,$(PLATFORMS)) $(addprefix release-,$(PLATFORMS))
build-all: $(addprefix build-,$(PLATFORMS))
release-all: $(addprefix release-,$(PLATFORMS))

# Windows 使用 .exe；默认禁用 CGO，因此交叉构建不需要额外 C 工具链。
define build_platform
	@target="$*"; target_os="$${target%-*}"; target_arch="$${target#*-}"; \
	ext=""; if [ "$$target_os" = windows ]; then ext=".exe"; fi; \
	output="$(BUILD_DIR)/$$target/$(BINARY_NAME)$$ext"; \
	mkdir -p "$(BUILD_DIR)/$$target" && \
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$$target_os GOARCH=$$target_arch $(GO) build $(1) $(LDFLAGS) -o "$$output" ./main.go
endef

$(addprefix build-,$(PLATFORMS)): build-%:
	$(call build_platform,)

$(addprefix release-,$(PLATFORMS)): release-%:
	$(call build_platform,-trimpath)
