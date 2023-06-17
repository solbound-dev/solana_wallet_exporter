ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest
LABEL maintainer="Bartol Deak <b@bdeak.net>"

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/solana_wallet_exporter /bin/solana_wallet_exporter

USER nobody
EXPOSE 18899
ENTRYPOINT [ "/bin/solana_wallet_exporter" ]
