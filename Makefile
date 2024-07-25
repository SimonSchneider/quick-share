# Description: Makefile for building the docker image
DOCKER_TAG=latest
DOCKER_IMAGE=ghcr.io/simonschneider/quick-share
DOCKER=docker

build_docker:
	$(DOCKER) build -t $(DOCKER_IMAGE) .

push:
	$(DOCKER) tag $(DOCKER_IMAGE) $(DOCKER_IMAGE):latest
	$(DOCKER) tag $(DOCKER_IMAGE) $(DOCKER_IMAGE):$(DOCKER_TAG)
	$(DOCKER) push $(DOCKER_IMAGE):latest
	$(DOCKER) push $(DOCKER_IMAGE):$(DOCKER_TAG)

build: build_docker push
