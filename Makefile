# ENVIRONMENT
PWD      := $(shell pwd)
MYSELF   := $(shell id -u)
MY_GROUP := $(shell id -g)

# PATHS
THIS                   := github.com/verygoodsoftwarenotvirus/typewizard
ARTIFACTS_DIR          := artifacts
COVERAGE_OUT           := $(ARTIFACTS_DIR)/coverage.out

# COMPUTED
TOTAL_PACKAGE_LIST    := `go list $(THIS)/...`

# CONTAINER VERSIONS
LINTER_IMAGE           := golangci/golangci-lint:v1.61.0

# COMMANDS
GO_FORMAT             := gofmt -s -w
GO_TEST               := CGO_ENABLED=1 go test -shuffle=on -race -vet=all
CONTAINER_RUNNER      := docker
RUN_CONTAINER         := $(CONTAINER_RUNNER) run --rm --volume $(PWD):$(PWD) --workdir=$(PWD)
RUN_CONTAINER_AS_USER := $(CONTAINER_RUNNER) run --rm --volume $(PWD):$(PWD) --workdir=$(PWD) --user $(MYSELF):$(MY_GROUP)
LINTER                := $(RUN_CONTAINER) $(LINTER_IMAGE) golangci-lint

## PREREQUISITES

.PHONY: ensure_fieldalignment_installed
ensure_fieldalignment_installed:
ifeq (, $(shell which fieldalignment))
	$(shell go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@v0.29.0)
endif

.PHONY: ensure_tagalign_installed
ensure_tagalign_installed:
ifeq (, $(shell which tagalign))
	$(shell go install github.com/4meepo/tagalign/cmd/tagalign@v1.4.1)
endif

.PHONY: ensure_gci_installed
ensure_gci_installed:
ifeq (, $(shell which gci))
	$(shell go install github.com/daixiang0/gci@v0.13.5)
endif

.PHONY: ensure_goimports_installed
ensure_goimports_installed:
ifeq (, $(shell which goimports))
	$(shell go install golang.org/x/tools/cmd/goimports@v0.29.0)
endif

.PHONY: clean_vendor
clean_vendor:
	rm -rf vendor go.sum

vendor:
	if [ ! -f go.mod ]; then go mod init; fi
	go mod tidy
	go mod vendor

.PHONY: revendor
revendor: clean_vendor vendor

## FORMATTING

.PHONY: format_golang
format_golang: format_imports ensure_fieldalignment_installed ensure_tagalign_installed
	@until fieldalignment -fix ./...; do true; done > /dev/null
	@until tagalign -fix -sort -order "env,envDefault,envPrefix,json,mapstructure,toml,yaml" ./...; do true; done > /dev/null
	for file in `find $(PWD) -type f -not -path '*/vendor/*' -name "*.go"`; do $(GO_FORMAT) $$file; done

.PHONY: format_imports
format_imports: ensure_gci_installed
	gci write --section standard --section "prefix($(THIS))" --section "prefix($(dir $(THIS)))" --section default --custom-order `find $(PWD) -type f -not -path '*/vendor/*' -name "*.go"`

.PHONY: format
format: format_golang

.PHONY: fmt
fmt: format

.PHONY: goimports
goimports: ensure_goimports_installed
	goimports -w .

## LINTING

.PHONY: golang_lint
golang_lint:
	@$(CONTAINER_RUNNER) pull --quiet $(LINTER_IMAGE)
	$(LINTER) run --config=.golangci.yml --timeout 15m ./...

.PHONY: lint
lint: golang_lint

.PHONY: clean_coverage
clean_coverage:
	@rm --force $(COVERAGE_OUT) profile.out;

.PHONY: coverage
coverage: clean_coverage $(ARTIFACTS_DIR)
	@$(GO_TEST) -coverprofile=$(COVERAGE_OUT) -covermode=atomic $(TOTAL_PACKAGE_LIST) > /dev/null
	@go tool cover -func=$(ARTIFACTS_DIR)/coverage.out | grep 'total:' | xargs | awk '{ print "COVERAGE: " $$3 }'

## EXECUTION

.PHONY: build
build:
	go build $(TOTAL_PACKAGE_LIST)

.PHONY: test
test: vendor build
	$(GO_TEST) -failfast $(TOTAL_PACKAGE_LIST)
