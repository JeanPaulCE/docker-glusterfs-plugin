ARG UBUNTU_VERSION=22.04
FROM --platform=${TARGETPLATFORM:-linux/amd64} ubuntu:${UBUNTU_VERSION} as base

ENV DEBIAN_FRONTEND="noninteractive"

# Install required packages and clean up in a single layer
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    glusterfs-client \
    curl \
    rsyslog \
    tini && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf "/tmp/*" "/root/.cache" /var/log/lastlog && \
    mkdir -p /var/lib/glusterd /etc/glusterfs && \
    touch /etc/glusterfs/logger.conf

COPY rsyslog.conf /etc/rsyslog.conf

ARG GO_VERSION=1.21.6
FROM --platform=${TARGETPLATFORM:-linux/amd64} golang:${GO_VERSION}-alpine as dev

RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/glusterfs-volume-plugin ./cmd/main.go

FROM base

COPY --from=dev /go/bin/glusterfs-volume-plugin /

ENTRYPOINT ["/usr/bin/tini", "--", "/glusterfs-volume-plugin"]
