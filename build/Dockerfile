FROM golang:1.17-alpine3.15 AS builder

RUN addgroup -g 1001 wowpow && \
    adduser -S -u 1001 -G wowpow wowpow

RUN \
    apk add --update --no-cache \
        make \
        bash \
        git \
        openssh-client \
        build-base \
        musl-dev \
        curl \
        wget \
        tzdata

ADD . /src/app

WORKDIR /src/app

ARG TARGETPLATFORM
ARG APPLICATION=wowserver
RUN \
    make gen && \
    CGO_ENABLED=0 go build  -ldflags="-d -s -w" \
    -o "./bin/app" \
        cmd/${APPLICATION}/main.go

FROM alpine
WORKDIR /

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /src/app/bin/app /app

USER wowpow:wowpow

ENTRYPOINT ["/app"]
