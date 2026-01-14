IMAGE_NAME := tetra_bot
TAG := latest
REGISTRY ?= "ghcr.io/piterpentester"

# Helper to check if a command exists
HAS_GOLANGCI := $(shell command -v golangci-lint 2> /dev/null)

.PHONY: all build image push clean deploy test lint fmt help k3s-import

help: ## Show this help message
	@echo "Tetra (Time to Restart) - Internet Monitor Bot"
	@echo
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: build ## Build binary (default)

build: ## Build the Go binary for the local architecture
	go build -o tetra ./cmd/tetra

test: ## Run unit tests
	go test -v ./...

lint: ## Run linters (vet or golangci-lint if installed)
ifdef HAS_GOLANGCI
	golangci-lint run
else
	go vet ./...
	@echo "Tip: Install golangci-lint for better checks"
endif

fmt: ## Format code
	go fmt ./...

image: ## Build Docker image
	docker build -t $(IMAGE_NAME):$(TAG) .

k3s-import: image ## Build image and import into k3s (for local dev on Pi)
	sudo k3s ctr images import $(IMAGE_NAME).tar || \
	docker save $(IMAGE_NAME):$(TAG) | sudo k3s ctr images import -

push: ## Push Docker image to registry (requires REGISTRY env)
	@if [ -z "$(REGISTRY)" ]; then echo "Error: REGISTRY not set. Use make push REGISTRY=ghcr.io/user"; exit 1; fi
	docker tag $(IMAGE_NAME):$(TAG) $(REGISTRY)/$(IMAGE_NAME):$(TAG)
	docker push $(REGISTRY)/$(IMAGE_NAME):$(TAG)

deploy: ## Deploy to Kubernetes (Apply k8s/ yaml files)
	kubectl apply -f k8s/configmap.yaml
	kubectl apply -f k8s/secrets.yaml
	kubectl apply -f k8s/deployment.yaml

clean: ## Remove built binary
	rm -f tetra
