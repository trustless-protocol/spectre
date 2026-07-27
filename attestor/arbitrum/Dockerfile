# syntax=docker/dockerfile:1.7

ARG NITRO_IMAGE=offchainlabs/nitro-node:v3.11.2-3599aca
FROM ${NITRO_IMAGE} AS nitro-runtime

FROM golang:1.25.7-bookworm AS attestor-builder

WORKDIR /src/arbitrum

COPY arbitrum/go.mod arbitrum/go.sum ./
COPY types/go.mod types/go.sum ../types/
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY types ../types
COPY arbitrum/*.go ./
COPY arbitrum/cmd ./cmd
COPY arbitrum/server ./server
COPY arbitrum/bold ./bold
COPY arbitrum/legacy ./legacy

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" \
    -o /out/attestor ./cmd

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install --no-install-recommends -y ca-certificates jq libstdc++6 \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd --gid 10001 attestor \
    && useradd --uid 10001 --gid attestor --create-home attestor \
    && mkdir -p /config /run/fast-ibc /usr/local/share/fast-ibc \
        /var/lib/fast-ibc/nitro /home/user/nitro-legacy /home/user/target \
    && chown -R attestor:attestor /run/fast-ibc /var/lib/fast-ibc

COPY --from=attestor-builder /out/attestor /usr/local/bin/attestor
COPY --from=nitro-runtime --chmod=0755 /usr/local/bin/nitro /usr/local/bin/nitro
COPY --from=nitro-runtime /home/user/nitro-legacy/machines /home/user/nitro-legacy/machines
COPY --from=nitro-runtime /home/user/target/machines /home/user/target/machines
COPY --chmod=0755 arbitrum/docker/entrypoint.sh /usr/local/bin/attestor-entrypoint

RUN sha256sum /usr/local/bin/nitro \
    | awk '{print $1}' > /usr/local/share/fast-ibc/nitro.sha256
RUN /usr/local/bin/nitro --help >/dev/null 2>&1; \
    status=$?; \
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]

USER attestor
WORKDIR /config

ENV ATTESTOR_CONFIG=/config/config.json \
    ATTESTOR_GRPC_LISTEN_ADDRESS=0.0.0.0:50051 \
    NITRO_BINARY_PATH=/usr/local/bin/nitro \
    NITRO_BINARY_SHA256_FILE=/usr/local/share/fast-ibc/nitro.sha256 \
    NITRO_ALLOWED_WASM_MODULE_ROOTS=/home/user/nitro-legacy/machines,/home/user/target/machines \
    NITRO_DATA_DIR=/var/lib/fast-ibc/nitro \
    NITRO_IPC_PATH=/run/fast-ibc/nitro.ipc \
    ATTESTOR_STATE_PATH=/var/lib/fast-ibc/nitro/attested-roots.json

VOLUME ["/var/lib/fast-ibc/nitro"]
EXPOSE 50051

ENTRYPOINT ["/usr/local/bin/attestor-entrypoint"]
CMD ["start"]
