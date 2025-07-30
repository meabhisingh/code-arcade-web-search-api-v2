# syntax=docker/dockerfile:1

################################################################################
# Build stage
ARG GO_VERSION=1.24.3
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Set architecture and build the binary
ARG TARGETARCH
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server ./cmd

################################################################################
# Final runtime stage
FROM alpine:3.22 AS final

# Install runtime dependencies
RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
    ca-certificates \
    tzdata \
    curl \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

# Create non-root user
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

# Copy built binary
COPY --from=build /bin/server /bin/
COPY --from=build /src/public /public

# Expose application port
EXPOSE 3000

# Add a healthcheck using /test route
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD curl -sf http://localhost:3000/test | grep -q healthy || exit 1

# Entrypoint
ENTRYPOINT [ "/bin/server" ]
