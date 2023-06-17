DOCKER_REPO ?= ghcr.io/solbound-dev
DOCKER_IMAGE_NAME ?= solana-wallet-exporter
DOCKER_ARCHS ?= amd64 armv7 arm64

include Makefile.common

.PHONY: crossbuild
crossbuild: promu
	@echo ">> cross-building binaries"
	$(PROMU) crossbuild

.PHONY: tarballs
tarballs: promu
	@echo ">> building release tarballs"
	$(PROMU) crossbuild tarballs

.PHONY: docker-push-latest $(TAG_DOCKER_ARCHS)
docker-push-latest: $(TAG_DOCKER_ARCHS)
$(TAG_DOCKER_ARCHS): common-docker-tag-latest-%:
	docker push "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME)-linux-$*:latest"
	docker push "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME)-linux-$*:v$(DOCKER_MAJOR_VERSION_TAG)"

.PHONY: common-docker-manifest-latest
common-docker-manifest-latest:
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest create -a "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):latest" $(foreach ARCH,$(DOCKER_ARCHS),$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME)-linux-$(ARCH):latest)
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest push "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):latest"
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest create -a "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):v$(DOCKER_MAJOR_VERSION_TAG)" $(foreach ARCH,$(DOCKER_ARCHS),$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME)-linux-$(ARCH):v$(DOCKER_MAJOR_VERSION_TAG))
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest push "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):v$(DOCKER_MAJOR_VERSION_TAG)"
