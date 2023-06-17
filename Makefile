DOCKER_REPO ?= ghcr.io/solbound-dev
DOCKER_IMAGE_NAME ?= solana_wallet_exporter
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