# Stage 1: Build the Go Backend
FROM golang:1.22-alpine AS backend-builder

WORKDIR /app
COPY . .

# Add build arguments for debug mode and versioning
ARG DEBUG=false
ARG GIN_MODE=release
ENV GIN_MODE=${GIN_MODE}
ARG VERSION
ENV VERSION=${VERSION}

# Set build tags based on the DEBUG flag and include versioning information
RUN if [ "$DEBUG" = "true" ]; then \
      go mod download && go build -o server -tags debug -ldflags="-X main.version=DEBUG" ./cmd/server/main.go; \
    else \
      go mod download && go build -o server -ldflags="-X main.version=${VERSION}" ./cmd/server/main.go; \
    fi

# Stage 2: Combine Backend and Frontend into a Single Image
# TODO: replace this with minimal image
FROM nginx:alpine

# Add OCI Image Spec labels
ARG GIT_COMMIT_SHA
ARG GIT_VERSION

# TODO: replace these labels
LABEL org.opencontainers.image.title="event-go" \
      org.opencontainers.image.description="Event Go, Broadcast server sent events by sending HTTP calls." \
      org.opencontainers.image.source="https://github.com/heissonwillen/event-go" \
      org.opencontainers.image.url="https://github.com/heissonwillen/event-go" \
      org.opencontainers.image.documentation="https://github.com/heissonwillen/event-go" \
      org.opencontainers.image.licenses="MIT"

# Copy the Go server binary
COPY --from=backend-builder /app/server /app/server

# Create the storage directory for the backend
RUN mkdir -p /app/storage

# Expose the ports for Nginx and the Go server
EXPOSE 80 8080

# TODO: run Go binary directly
# Run both the Go server and Nginx using a simple script
CMD ["/bin/sh", "-c", "/app/server & nginx -g 'daemon off;'"]
