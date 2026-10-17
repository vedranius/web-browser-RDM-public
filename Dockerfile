# Web Remote Manager PRO — container image
#   docker build -t wrm-pro .
#   docker run -d -p 8080:8080 -v wrm-data:/data wrm-pro
# Data (database, encryption key, recordings, self-signed certificate) lives in /data.

FROM golang:1.27-alpine AS build
WORKDIR /src
COPY remote-manager/go.mod remote-manager/go.sum ./
RUN go mod download
COPY remote-manager/ ./
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.AppVersion=${VERSION}" -o /out/wrm .

FROM node:22-alpine3.22
# Node runs the Bitwarden CLI. npm, corepack and yarn are only needed to install it, so they
# are removed afterwards (their bundled modules would be outdated code in the image).
LABEL org.opencontainers.image.title="Web Remote Manager PRO" \
      org.opencontainers.image.source="https://github.com/vedranius/web-browser-RDM-public" \
      org.opencontainers.image.licenses="PolyForm-Noncommercial-1.0.0 OR PolyForm-Internal-Use-1.0.0"
RUN apk add --no-cache ca-certificates tzdata python3 py3-pip git openssh-client rsync \
 && npm install --global --omit=dev @bitwarden/cli@2026.9.1 \
 && npm cache clean --force \
 && rm -rf /usr/local/lib/node_modules/npm /usr/local/lib/node_modules/corepack /opt/yarn-* \
      /usr/local/bin/npm /usr/local/bin/npx /usr/local/bin/corepack /usr/local/bin/yarn /usr/local/bin/yarnpkg \
 && adduser -D -u 10001 wrm \
 && mkdir -p /data && chown wrm:wrm /data && chmod 700 /data
COPY --from=build /out/wrm /usr/local/bin/wrm
USER wrm
WORKDIR /data
ENV DB_PATH=/data/remote_manager.db \
    PORT=8080
EXPOSE 8080 3478/tcp 3478/udp
VOLUME ["/data"]
HEALTHCHECK --interval=30s --timeout=6s --start-period=10s --retries=3 CMD ["/usr/local/bin/wrm", "-healthcheck"]
ENTRYPOINT ["/usr/local/bin/wrm"]
