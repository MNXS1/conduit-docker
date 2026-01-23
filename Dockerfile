FROM --platform=linux/amd64 golang:1.24-alpine AS builder

RUN apk add --no-cache git make gcc musl-dev

WORKDIR /build

ARG PSIPHON_REPO_PATH=../psiphon-tunnel-core

COPY ${PSIPHON_REPO_PATH} /build/psiphon-tunnel-core/

WORKDIR /build/psiphon-tunnel-core

RUN go mod download

COPY main.go /build/psiphon-tunnel-core/main.go

WORKDIR /build/psiphon-tunnel-core

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -tags PSIPHON_ENABLE_INPROXY \
    -ldflags "-s -w" \
    -o /build/inproxy-node \
    main.go

FROM --platform=linux/amd64 alpine:latest

RUN apk add --no-cache ca-certificates && \
    mkdir -p /data && \
    addgroup -g 1000 inproxy && \
    adduser -D -u 1000 -G inproxy inproxy && \
    chown -R inproxy:inproxy /data

COPY --from=builder /build/inproxy-node /usr/local/bin/inproxy-node

WORKDIR /data

USER inproxy

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD pgrep -f inproxy-node > /dev/null || exit 1

ENTRYPOINT ["/usr/local/bin/inproxy-node"]

CMD ["-dataRootDirectory", "/data"]
