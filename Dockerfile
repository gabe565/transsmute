# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.24.0-alpine AS go-builder
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
    && go build -ldflags='-w -s' -trimpath -tags grpcnotrace


FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=go-builder /app/transsmute ./
ENV TRANSSMUTE_ADDRESS=":80"
CMD ["/transsmute"]
