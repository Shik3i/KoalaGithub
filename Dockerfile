# syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM node:24.18.0-alpine3.24 AS frontend-builder
WORKDIR /app

ARG APP_VERSION=dev
ENV PUBLIC_APP_VERSION=$APP_VERSION

COPY package.json package-lock.json svelte.config.js tsconfig.json vite.config.ts ./
RUN --mount=type=cache,target=/root/.npm,sharing=locked npm ci --ignore-scripts
COPY src ./src
COPY static ./static
RUN npm run prepare && npm run build

FROM --platform=$BUILDPLATFORM golang:1.27rc2-alpine3.24 AS backend-builder
WORKDIR /app

ARG TARGETOS
ARG TARGETARCH
ARG APP_VERSION=dev

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked go mod download
COPY internal ./internal
COPY main.go ./
COPY src/lib/data/visualizers.json ./src/lib/data/visualizers.json
COPY --from=frontend-builder /app/www ./www
RUN --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
	CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
	go build -trimpath -ldflags="-s -w -X main.version=${APP_VERSION}" -o /out/koala-github .

FROM alpine:3.24.1

RUN apk add --no-cache ca-certificates su-exec tzdata \
	&& addgroup -S -g 10001 koala \
	&& adduser -S -D -H -u 10001 -G koala koala \
	&& install -d -o koala -g koala /data

WORKDIR /app
COPY --from=backend-builder --chown=koala:koala /out/koala-github /app/koala-github
COPY --chmod=755 docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

ENV PORT=8080 \
	DB_PATH=/data/koalagithub.db

VOLUME ["/data"]
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
	CMD wget -qO- http://127.0.0.1:8080/api/health >/dev/null || exit 1

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["/app/koala-github"]
