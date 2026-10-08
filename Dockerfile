# Stage 1: Static executable compiler core for Go binary
FROM golang:alpine AS builder

WORKDIR /app
# hadolint ignore=DL3016,DL3018
RUN apk add --no-cache git curl tar nodejs npm \
    && npm install -g typescript esbuild \
    && curl -L https://github.com/sass/dart-sass/releases/download/1.105.0/dart-sass-1.105.0-linux-x64-musl.tar.gz -o /tmp/sass.tar.gz \
    && tar -xzf /tmp/sass.tar.gz -C /usr/local/share/ \
    && ln -s /usr/local/share/dart-sass/sass /usr/local/bin/sass \
    && rm /tmp/sass.tar.gz

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p html/css html/js\
    && sass html/sass/main.scss html/css/layout.css --style=compressed \
    && tsc --noEmit \
    && esbuild html/ts/main.ts --bundle \
    # --minify \
    --target=es2022 --format=iife --outfile=html/js/main.min.js \
    && rm -R html/sass html/ts \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -extldflags '-static'" -o audiobook-ingester .

# Stage 2: High-Performance, Minimal Production Runtime Container
# hadolint ignore=DL3007
FROM alpine:latest

# hadolint ignore=DL3018
RUN apk add --no-cache ffmpeg ca-certificates wget gcompat tzdata \
    && echo "hosts: files dns" > /etc/nsswitch.conf \
    && wget --progress=dot:giga https://github.com/djdembeck/m4b-merge/releases/download/v1.0.0/m4b-merge-linux -O /usr/local/bin/m4b-merge \
    && chmod +x /usr/local/bin/m4b-merge

COPY --from=builder /app/audiobook-ingester /audiobook-ingester

ENTRYPOINT ["/audiobook-ingester"]