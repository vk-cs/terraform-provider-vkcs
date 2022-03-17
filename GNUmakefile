TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=vkcs

default: build

build: fmtcheck
	go install

build_darwin: fmtcheck
	GOOS=darwin CGO_ENABLED=0 go build -o terraform-provider-vkcs_darwin

build_linux: fmtcheck
	GOOS=linux CGO_ENABLED=0 go build -o terraform-provider-vkcs_linux

build_windows: fmtcheck
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

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1
	golangci-lint run ./...

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile website website-test lint

deploy_linux: build_linux
	cp terraform-provider-vkcs_linux ~/.terraform.d/plugins/hub.vkcs.mail.ru/repository/vkcs/0.6.0/linux_amd64/terraform-provider-vkcs_0.6.0

