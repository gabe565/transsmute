# syntax=docker/dockerfile:1.2
ARG GO_VERSION=1.18

FROM --platform=$BUILDPLATFORM golang:$GO_VERSION-alpine as go-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY . .

# Set Golang build envs based on Docker platform string
ARG TARGETPLATFORM
RUN --mount=type=cache,id="transsmute-$TARGETPLATFORM",target=/root/.cache \
    set -x \
    && case "$TARGETPLATFORM" in \
        'linux/amd64') export GOARCH=amd64 ;; \
        'linux/arm/v6') export GOARCH=arm GOARM=6 ;; \
        'linux/arm/v7') export GOARCH=arm GOARM=7 ;; \
        'linux/arm64') export GOARCH=arm64 ;; \
        *) echo "Unsupported target: $TARGETPLATFORM" && exit 1 ;; \
    esac \
    && go build -ldflags='-w -s'


FROM alpine
LABEL org.opencontainers.image.source="https://github.com/gabe565/transsmute"
WORKDIR /app

RUN apk add --no-cache tzdata

COPY --from=go-builder /app/transsmute ./

ARG USERNAME=transsmute
ARG UID=1000
ARG GID=$UID
RUN addgroup -g "$GID" "$USERNAME" \
    && adduser -S -u "$UID" -G "$USERNAME" "$USERNAME"
USER $UID

ENV TRANSSMUTE_ADDRESS=":80"
CMD ["./transsmute"]
