# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.6.1 AS xx

FROM --platform=$BUILDPLATFORM golang:1.24.1-alpine AS go-builder
WORKDIR /app

COPY --from=xx / /

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY . .

ARG TARGETPLATFORM
RUN --mount=type=cache,id="transsmute-$TARGETPLATFORM",target=/root/.cache \
    CGO_ENABLED=0 xx-go build -ldflags='-w -s' -trimpath -tags grpcnotrace


FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=go-builder /app/transsmute ./
ENV TRANSSMUTE_ADDRESS=":80"
CMD ["/transsmute"]
