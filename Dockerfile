# SQLites needs build-essential and libsqlite3, that's why Debian is used here
FROM golang:1.22-bookworm AS build

WORKDIR /src
COPY . /src/

ARG DEBUG=false
ARG GIN_MODE=release
ENV GIN_MODE=${GIN_MODE}
ARG VERSION
ENV VERSION=${VERSION}

# TODO: different build for debug env
RUN CGO_ENABLED=1 GOOS=linux go mod download && go build -o /bin/event-go -a -ldflags '-linkmode external -extldflags "-static"' ./cmd/server/main.go

FROM scratch

ARG GIT_COMMIT_SHA
ARG GIT_VERSION

LABEL org.opencontainers.image.title="event-go" \
      org.opencontainers.image.description="Event Go, Broadcast server sent events by sending HTTP calls." \
      org.opencontainers.image.source="https://github.com/heissonwillen/event-go" \
      org.opencontainers.image.url="https://github.com/heissonwillen/event-go" \
      org.opencontainers.image.documentation="https://github.com/heissonwillen/event-go" \
      org.opencontainers.image.licenses="MIT"

COPY --from=build /bin/event-go /bin/event-go
CMD ["/bin/event-go"]
