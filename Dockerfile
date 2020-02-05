FROM golang:1.13 as builder
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /flatsearch

FROM alpine:3.11
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /flatsearch /
COPY --from=builder /src/config.json /config.json
USER appuser:appuser
ENTRYPOINT ["/flatsearch", "-config=/config.json"]
