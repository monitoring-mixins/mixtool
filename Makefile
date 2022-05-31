all: check-license build generate test

GITHUB_URL=github.com/monitoring-mixins/mixtool
GOOS?=$(shell uname -s | tr A-Z a-z)
GOARCH?=$(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m)))
OUT_DIR=_output
BIN?=mixtool
VERSION?=$(shell cat VERSION)
PKGS=$(shell go list ./... | grep -v /vendor/)
BUILDFLAGS?=-ldflags="-s -w -X main.version=$(VERSION)" -gcflags="-trimpath=$(GOPATH)" -asmflags="-trimpath=$(GOPATH)"

check-license:
	@echo ">> checking license headers"
	@./scripts/check_license.sh

crossbuild:
	@GOOS=linux ARCH=amd64 $(MAKE) -s build

build:
	@$(eval OUTPUT=$(OUT_DIR)/$(GOOS)/$(GOARCH)/$(BIN))
	@echo ">> building for $(GOOS)/$(GOARCH) to $(OUTPUT)"
	@mkdir -p $(OUT_DIR)/$(GOOS)/$(GOARCH)
	@CGO_ENABLED=0 go build --installsuffix cgo -o $(OUTPUT) $(BUILDFLAGS) $(GITHUB_URL)/cmd/$(BIN)

install: build
	@$(eval OUTPUT=$(OUT_DIR)/$(GOOS)/$(GOARCH)/$(BIN))
	@echo ">> copying $(BIN) into $(GOPATH)/bin/$(BIN)"
	@cp $(OUTPUT) $(GOPATH)/bin/$(BIN)

test:
	@echo ">> running all tests"
	@go test $(PKGS)

generate: embedmd
	@echo ">> generating docs"
	@./scripts/generate-help-txt.sh
	@embedmd -w `find ./ -path ./vendor -prune -o -name "*.md" -print`

embedmd:
	@go install github.com/campoy/embedmd@v1.0.0

.PHONY: all check-license crossbuild build install test generate embedmd
