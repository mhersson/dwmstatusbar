##
# DWM status bar Makefile
#
shell=bash

GO_FUMPT_EXISTS := $(shell command -v gofumpt)
GOLANGCI_LINT_EXISTS := $(shell command -v golangci-lint)
UPX_EXISTS := $(shell command -v upx)

LDFLAGS="-s -w"

OUT_BIN = dwmstatusbar

# The directory where the binary will be installed
INSTALL_DIR = /usr/local/bin

all: compress

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

test: ## run tests
	@go test -v ./... --coverprofile=cover.out

vet:
ifdef GOLANGCI_LINT_EXISTS
	@golangci-lint run
else
	@echo "Using go vet to check the code..."
	@go vet ./...
endif

fmt:

ifdef GO_FUMPT_EXISTS
	@gofumpt -w .
else
	@echo "Using go fmt to format the code..."
	@go fmt ./...
endif

build: fmt vet test ## build the code
	@go build -ldflags $(LDFLAGS) -o bin/${OUT_BIN}

compress: build ## compress the binary
ifdef UPX_EXISTS
	@upx -q bin/${OUT_BIN}
endif

install: ## install the binary
	@cp bin/$(OUT_BIN) $(INSTALL_DIR)

uninstall: ## uninstall the binary
	@rm -f $(INSTALL_DIR)/$(OUT_BIN)

clean: ## clean up the project
	@rm -rf bin
