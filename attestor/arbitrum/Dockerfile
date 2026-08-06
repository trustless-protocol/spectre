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
    && mkdir -p /config /run/fast-ibc /var/lib/fast-ibc/attestor \
    && chown -R attestor:attestor /run/fast-ibc /var/lib/fast-ibc/attestor

COPY --from=attestor-builder /out/attestor /usr/local/bin/attestor
COPY --chmod=0755 arbitrum/docker/entrypoint.sh /usr/local/bin/attestor-entrypoint

USER attestor
WORKDIR /config

ENV ATTESTOR_CONFIG=/config/config.json \
    ATTESTOR_GRPC_LISTEN_ADDRESS=0.0.0.0:50051 \
    ATTESTOR_STATE_PATH=/var/lib/fast-ibc/attestor/attested-roots.json

VOLUME ["/var/lib/fast-ibc/attestor"]
EXPOSE 50051

ENTRYPOINT ["/usr/local/bin/attestor-entrypoint"]
CMD ["start"]
