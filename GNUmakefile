TEST        ?= $$(go list ./... |grep -v 'vendor')
GOFMT_FILES ?= $$(find . -name '*.go' |grep -v vendor)
PKG_NAME     = vkcs
GO          ?= go


LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

GOLANGCI_LINT         ?= $(LOCALBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v2.7.2
TF_PLUGIN_GEN         ?= $(LOCALBIN)/tfplugingen-framework
TF_PLUGIN_GEN_VERSION ?= v0.4.1

define go-install-tool
@[ -f $(LOCALBIN)/$(1) ] || { \
set -e ;\
echo "Installing $(1)@$(3)" ;\
GOBIN=$(LOCALBIN) $(GO) install $(2)@$(3) ;\
}
endef

default: build

golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

tfplugingen-framework: $(TF_PLUGIN_GEN)
$(TF_PLUGIN_GEN): $(LOCALBIN)
	$(call go-install-tool,$(TF_PLUGIN_GEN),github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework,$(TF_PLUGIN_GEN_VERSION))

generate: tfplugingen-framework
	PATH="$(LOCALBIN):$$PATH" go generate ./...

build: fmtcheck generate
	go install

build_darwin: fmtcheck generate
	GOOS=darwin CGO_ENABLED=0 go build -o terraform-provider-vkcs_darwin

build_linux: fmtcheck generate
	GOOS=linux CGO_ENABLED=0 go build -o terraform-provider-vkcs_linux

build_windows: fmtcheck generate
	GOOS=windows CGO_ENABLED=0 go build -o terraform-provider-vkcs_windows

test: fmtcheck
	go test $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -cover -timeout=30s -parallel=4

testacc_compute: fmtcheck
	TF_ACC=1 go test -run=TestAccCompute $(TEST) -v $(TESTARGS) -timeout 120m

testacc_image: fmtcheck
	TF_ACC=1 go test -run=TestAccImagesImage $(TEST) -v $(TESTARGS) -timeout 120m

testacc_keymanager: fmtcheck
	TF_ACC=1 go test -run=TestAccKeyManager $(TEST) -v $(TESTARGS) -timeout 120m

testacc_blockstorage: fmtcheck
	TF_ACC=1 go test -run=TestAccBlockStorage $(TEST) -v $(TESTARGS) -timeout 120m

testacc_lb: fmtcheck
	TF_ACC=1 go test -run=TestAccLB $(TEST) -v $(TESTARGS) -timeout 120m

testacc_vpnaas: fmtcheck
	TF_ACC=1 go test -run=TestAccVPNaaS $(TEST) -v $(TESTARGS) -timeout 120m

testacc_sfs: fmtcheck
	TF_ACC=1 go test -run=TestAccSFS $(TEST) -v $(TESTARGS) -timeout 120m

testacc_networking: fmtcheck
	TF_ACC=1 go test -run=TestAccNetworking $(TEST) -v $(TESTARGS) -timeout 120m

testacc_dbaas: fmtcheck
	TF_ACC=1 go test -run=TestAccDatabase $(TEST) -v $(TESTARGS) -timeout 120m

testmock_k8saas: fmtcheck
	TF_ACC=1 TF_ACC_MOCK_MCS=1 go test $(TEST) -run=TestMockAcc $(TESTARGS) -timeout 120m

testacc_k8saas: fmtcheck
	TF_ACC=1 go test -run=TestAccKubernetes $(TEST) $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

lint: golangci-lint
	$(GOLANGCI_LINT) run ./...

update_release_schema:
	go run helpers/schema-api/main.go -export .release/provider-schema.json

tflint_fix_examples:
	tflint --chdir=examples --recursive -f compact --config="$(CURDIR)/.tflint.hcl" --fix

tffmt_examples:
	terraform fmt --recursive examples

check_examples:
	tflint --chdir=examples --recursive -f compact --config="$(CURDIR)/.tflint.hcl"
	terraform fmt --check --recursive examples

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile website website-test lint update_release_schema generate tfplugingen-framework
