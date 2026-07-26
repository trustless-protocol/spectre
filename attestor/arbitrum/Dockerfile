# syntax=docker/dockerfile:1.7

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

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" \
    -o /out/attestor ./cmd

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install --no-install-recommends -y ca-certificates jq libstdc++6 \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd --gid 10001 attestor \
    && useradd --uid 10001 --gid attestor --create-home attestor \
    && mkdir -p /config /run/fast-ibc /var/lib/fast-ibc/nitro \
    && chown -R attestor:attestor /run/fast-ibc /var/lib/fast-ibc

COPY --from=attestor-builder /out/attestor /usr/local/bin/attestor
COPY --chmod=0755 arbitrum/bin/nitro /usr/local/bin/nitro
COPY --chmod=0755 arbitrum/docker/entrypoint.sh /usr/local/bin/attestor-entrypoint

ARG NITRO_BINARY_SHA256
RUN test -n "$NITRO_BINARY_SHA256" \
    && printf '%s  %s\n' "$NITRO_BINARY_SHA256" /usr/local/bin/nitro \
    | sha256sum --check -
RUN /usr/local/bin/nitro --help >/dev/null 2>&1; \
    status=$?; \
    [ "$status" -eq 0 ] || [ "$status" -eq 1 ]

USER attestor
WORKDIR /config

ENV ATTESTOR_CONFIG=/config/config.json \
    ATTESTOR_GRPC_LISTEN_ADDRESS=0.0.0.0:50051 \
    NITRO_BINARY_PATH=/usr/local/bin/nitro \
    NITRO_BINARY_SHA256=$NITRO_BINARY_SHA256 \
    NITRO_DATA_DIR=/var/lib/fast-ibc/nitro \
    NITRO_IPC_PATH=/run/fast-ibc/nitro.ipc \
    ATTESTOR_STATE_PATH=/var/lib/fast-ibc/nitro/attested-roots.json

VOLUME ["/var/lib/fast-ibc/nitro"]
EXPOSE 50051

ENTRYPOINT ["/usr/local/bin/attestor-entrypoint"]
CMD ["start"]
