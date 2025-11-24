# Disable all the default make stuff
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

## Display a list of the documented make targets
.PHONY: help
help:
	@echo Documented Make targets:
	@perl -e 'undef $$/; while (<>) { while ($$_ =~ /## (.*?)(?:\n# .*)*\n.PHONY:\s+(\S+).*/mg) { printf "\033[36m%-30s\033[0m %s\n", $$2, $$1 } }' $(MAKEFILE_LIST) | sort

.PHONY: .FORCE
.FORCE:

build:
	go build ./cmd/score-radius/

test:
	go vet ./...    
	go test -v ./... -cover -race

build-container:
	docker build -t score-radius:local .

test-app: build
	./score-radius --version
	./score-radius init
	./score-radius generate score.yaml
	cat app.bicep

test-app-with-full: build
	./score-radius --version
	./score-radius init --no-sample
	./score-radius generate examples/score/score-full.yaml
	cat app.bicep

test-app-with-redis: build
	./score-radius --version
	./score-radius init --no-sample
	cp examples/provisioners/redis.provisioners.yaml .score-radius/
	./score-radius generate examples/score/score-redis.yaml -i ghcr.io/radius-project/samples/demo:latest
	cat app.bicep

test-app-with-podinfo: build
	./score-radius --version
	./score-radius init --no-sample
	cp examples/provisioners/redis.provisioners.yaml .score-radius/
	./score-radius generate examples/score/score-podinfo-with-redis.yaml -i ghcr.io/stefanprodan/podinfo
	cat app.bicep

test-container: build-container
	docker run --rm score-radius:local --version
	docker run --rm -v .:/score-radius score-radius:local init
	docker run --rm -v .:/score-radius score-radius:local generate score.yaml
	cat app.bicep

## Create a local Kind cluster.
.PHONY: setup-kind-cluster
setup-kind-cluster:
	./scripts/setup-kind-cluster.sh

## Deploy podinfo to Radius and Kind cluster.
.PHONY: deploy-podinfo-to-radius
deploy-podinfo-to-radius:
	rad workspace create kubernetes default
	rad group create default --workspace default
	rad env create default --group default
	./score-radius init --no-sample
	cp ./examples/provisioners/redis.provisioners.yaml .score-radius/
	cp ./examples/score/score-podinfo.yaml ./score.yaml
	./score-radius generate score.yaml -i ghcr.io/stefanprodan/podinfo:latest -o app.bicep
	cat app.bicep
	cp ./examples/bicepconfig.json ./
	rad deploy app.bicep --group default --application podinfo --environment default
	kubectl wait deployments/podinfo -n default-podinfo --for condition=Available --timeout=90s
	kubectl get all -n default-podinfo
