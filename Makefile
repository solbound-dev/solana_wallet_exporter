DOCKER_IMAGE_NAME ?= solana_wallet_exporter
DOCKER_REPO ?= ghcr.io/bdeak4

include Makefile.common

.PHONY: crossbuild
crossbuild: precheck style check_license lint yamllint unused promu
	@echo ">> cross-building packages"
	$(PROMU) crossbuild

.PHONY: package-tarballs
package-tarballs: crossbuild
	@echo ">> packaging tarballs"
	$(PROMU) crossbuild tarballs

.PHONY: package-docker-images
package-docker-images: crossbuild
	@echo ">> packaging docker images"
	docker buildx build --output "type=image,push=true" -f Dockerfile.cross \
	  --tag "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(SANITIZED_DOCKER_IMAGE_TAG)" \
		--platform=linux/amd64,linux/arm64,linux/arm/v7 $(DOCKERBUILD_CONTEXT)

.PHONY: package
package: package-tarballs package-docker-images