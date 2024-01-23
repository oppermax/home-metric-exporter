EXECUTABLE ?= exporter
SERVICE ?= home-metric-exporter
OWNER ?= oppermax
GO := CGO_ENABLED=0 go
DATE := $(shell date -u '+%FT%T%z')
GOLANGCI_LINT_VERSION := 1.50.1
LINT_TARGETS ?= $(shell $(GO) list -f '{{.Dir}}' ./... | sed -e"s|${CURDIR}/\(.*\)\$$|\1/...|g" )

DOCKER_IMAGE_TAG ?= latest
DOCKER_FILE ?= Dockerfile

SYSTEM       := $(shell uname -s | tr A-Z a-z)_$(shell uname -m | sed "s/x86_64/amd64/" | sed "s/armv6l/armv6/" )

GO_OS = linux
GO_ARCH = arm
GO_ARM = 6

PACKAGES = $(shell go list ./...)

.PHONY: all
all: build

.PHONY: clean
clean:
	$(GO) clean -i ./...
	rm -rf bin/

.PHONY: fmt
fmt:
	$(GO) fmt $(PACKAGES)

.PHONY: lint
lint: bin/golangci-lint-$(GOLANGCI_LINT_VERSION)
	$(GO_PREFIX) ./bin/golangci-lint-$(GOLANGCI_LINT_VERSION) run $(LINT_TARGETS)

.PHONY: create-golint-config
create-golint-config: .golangci.yml

bin/golangci-lint-$(GOLANGCI_LINT_VERSION):
	mkdir -p bin
	curl -sSLf \
		https://github.com/golangci/golangci-lint/releases/download/v$(GOLANGCI_LINT_VERSION)/golangci-lint-$(GOLANGCI_LINT_VERSION)-$(shell echo $(SYSTEM) | tr '_' '-').tar.gz \
		| tar xzOf - golangci-lint-$(GOLANGCI_LINT_VERSION)-$(shell echo $(SYSTEM) | tr '_' '-')/golangci-lint > bin/golangci-lint-$(GOLANGCI_LINT_VERSION) && chmod +x bin/golangci-lint-$(GOLANGCI_LINT_VERSION)



.PHONY: test
test:
	@for PKG in $(PACKAGES); do $(GO) test -cover $$PKG || exit 1; done;

.PHONY: build
build: clean
	env GOOS=$(GO_OS) GOARCH=$(GO_ARCH) GOARM=$(GO_ARM) CGO_ENABLED=0 $(GO) build -ldflags="-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" -o bin/$(SERVICE) ./exporter/

.PHONY: release
release:
	@which gox > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/mitchellh/gox; \
	fi
	CGO_ENABLED=0 gox -verbose -osarch '!darwin/386' -output="dist/$(EXECUTABLE)-{{.OS}}-{{.Arch}}"

.PHONY: docker-build
docker-build: build
	docker buildx build --platform linux/arm/v6 --no-cache -t $(OWNER)/$(SERVICE):$(DOCKER_IMAGE_TAG) -f $(DOCKER_FILE) .

.PHONY: docker-push
docker-push: docker-build
	docker push $(OWNER)/$(SERVICE):$(DOCKER_IMAGE_TAG)

.PHONY: deploy
deploy:
	docker pull $(OWNER)/$(SERVICE):$(DOCKER_IMAGE_TAG)
	docker compose up -d

.PHONY: run-locally
run-locally:
	$(GO) run exporter/main.go
