###############################################################################
#
#                   ❇❇❇ Welcome to the Makefile ❇❇❇
#                             Version: v0.1
#
# Include this in your project, but don't change it. Rather create a main
# Makegile and include this, adjust variables as needed. You probably want to
# set the $IMG variable for the docker image.
#
###############################################################################

COLUMNS               ?= $(shell stty size | awk '{print $$2}')
SETUP_STTY            ?= -stty columns $$(( $(COLUMNS) - 4 )) $(DEBUG)
UPPER                 ?= $(shell echo '$1' | tr '[:lower:]' '[:upper:]')
FORMAT                ?= sed "s/^/    /g"

BIN_DIR               ?= bin
BIN                   ?= $(BIN_DIR)/app
GOLANGCI_LINT         ?= $(shell pwd)/$(BIN_DIR)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.52.2
GOFUMPT               ?= $(shell pwd)/$(BIN_DIR)/gofumpt
GOFUMPT_VERSION       ?= v0.4.0
GOTESTSUM             ?= $(shell pwd)/$(BIN_DIR)/gotestsum
GOTESTSUM_VERSION     ?= 1.8.2
TEST_FLAGS            ?= ./...
TRIVY                 ?= $(shell pwd)/$(BIN_DIR)/trivy
TRIVY_VERSION         ?= v0.34.0

RM                    ?= rm
MKDIR                 ?= mkdir

GO                    ?= go
GOBUILD               ?= $(GO) build
GOGET                 ?= $(GO) get

GO_SRC                ?= $(shell find ./ -name '*.go')

OS                    ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
CGO_ENABLED           ?= 0
GOENV                 ?= CGO_ENABLED=$(CGO_ENABLED) GOOS=$(OS)
LDFLAGS               ?= -ldflags="-s -w"

ECHO                  ?= @echo "  "
ifeq ($(OS),linux)
	ECHO                ?= @echo -e "  "
endif

ARCH                  ?= $(shell uname -m | tr '[:upper:]' '[:lower:]')
ifeq ($(ARCH),x86_64)
	ARCH               = amd64
endif

V                     ?= 0
ifneq ($(V),1)
	Q                  = @
	DEBUG              = 2>/dev/null
endif

.PHONY: all
all: help

##@ General
.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; \
		printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} \
		/^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Testing
.PHONY: test
test: $(GOTESTSUM) ## Run unit tests
	$(Q)$(ECHO) "GO" $(call UPPER, $@)
	$(Q)$(GOTESTSUM) $(TEST_FLAGS) | $(FORMAT)

##@ Development
.PHONY: build
build: $(BIN) ## Build executables
$(BIN_DIR)/%: $(GO_SRC)
	$(Q)$(ECHO) "GO" $(call UPPER, $@)
	$(Q)$(MKDIR) -p $(BIN_DIR)
	$(Q)$(GOENV) $(GOBUILD) $(LDFLAGS) -o $@ main.go 2>&1 | $(FORMAT)

.PHONY: fmt
fmt: $(GOFUMPT) ## Run gofumpt
	$(Q)$(ECHO) "GO" $(call UPPER, $@)
	$(Q)$(GOFUMPT) -w . 2>&1 | $(FORMAT)

.PHONY: lint
lint: ${GOLANGCI_LINT} ## Run golangci-lint linter
	$(Q)$(ECHO) $(call UPPER, $@)
	$(Q)$(GOLANGCI_LINT) run --color always ./... 2>&1 | $(FORMAT)

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT) ## Run golangci-lint linter and perform fixes
	$(Q)$(ECHO) $(call UPPER, $@)
	$(Q)$(GOLANGCI_LINT) run --fix --color always ./... 2>&1 | $(FORMAT)

.PHONY: clean
clean: ## Remove bin dir
	$(Q)$(ECHO) $(call UPPER, $@)
	$(Q)$(RM) -rf bin/*

##@ Vulnability checking
.PHONY: check-dockerfile
check-dockerfile: $(TRIVY) ## Check Dockerfile
	$(Q)$(ECHO) $(call UPPER, $@)
	$(Q)$(SETUP_STTY)
	$(Q)$(TRIVY) --exit-code 1 config Dockerfile 2>&1 | $(FORMAT)

.PHONY: check-image
check-image: $(TRIVY) ## Check Image
	$(Q)$(ECHO) $(call UPPER, $@)
ifndef IMG
	$(error IMG variable is not set)
endif
	$(Q)$(SETUP_STTY)
	$(Q)$(TRIVY) --exit-code 1 image $(IMG) --ignore-unfixed --security-checks vuln --severity HIGH,CRITICAL 2>&1 | $(FORMAT)

##@ Publishing
.PHONY: build-docker
build-docker: $(BINS)## Build the docker image
	$(Q)$(ECHO) $(call UPPER, $@)
ifndef IMG
	$(error IMG variable is not set)
endif
	$(Q) DOCKER_BUILDKIT=1 docker build . -t ${IMG} | $(FORMAT)

$(GOTESTSUM):
	$(Q)$(ECHO) "GOTESTSUM"
	$(Q)$(MKDIR) -p $(BIN_DIR)
	$(Q) curl -sSfL https://github.com/gotestyourself/gotestsum/releases/download/v$(GOTESTSUM_VERSION)/gotestsum_$(GOTESTSUM_VERSION)_$(OS)_$(ARCH).tar.gz | tar zx -C bin/ gotestsum

$(GOFUMPT):
	$(Q)$(ECHO) "GOFUMPT"
	$(Q)$(MKDIR) -p $(BIN_DIR)
	$(Q) curl -sSfL -o$(GOFUMPT) https://github.com/mvdan/gofumpt/releases/download/$(GOFUMPT_VERSION)/gofumpt_$(GOFUMPT_VERSION)_$(OS)_$(ARCH)
	$(Q) chmod +x $(GOFUMPT)

$(GOLANGCI_LINT):
	$(Q)$(ECHO) "GOLANGCI_LINT"
	$(Q)$(MKDIR) -p $(BIN_DIR)
	$(Q) curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh $(DEBUG) | sh -s -- -b $(shell dirname $(GOLANGCI_LINT)) $(GOLANGCI_LINT_VERSION) $(DEBUG) | $(FORMAT)

$(TRIVY):
	$(Q)$(ECHO) "TRIVY"
	$(Q)$(MKDIR) -p $(BIN_DIR)
	$(Q) curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh $(DEBUG) | sh -s -- -b $(shell dirname $(TRIVY)) $(TRIVY_VERSION) $(DEBUG) | $(FORMAT)
