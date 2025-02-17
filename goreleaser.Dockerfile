FROM gcr.io/distroless/static:nonroot
LABEL org.opencontainers.image.source="https://github.com/gabe565/transsmute"
WORKDIR /
COPY transsmute /
ENV TRANSSMUTE_ADDRESS=":80"
CMD ["/transsmute"]
