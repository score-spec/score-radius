FROM --platform=$BUILDPLATFORM dhi.io/golang:1.26.4-alpine3.23-dev@sha256:3c62ff766e4118b29b7815a0732affa1ac9836a67f5c505d5468ae08f7092a05 AS builder

ARG VERSION
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

# Set the current working directory inside the container.
WORKDIR /go/src/github.com/score-spec/score-radius

# Copy just the module bits
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project and build it.
COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 \
    go build -ldflags="-s -w \
        -X github.com/score-spec/score-radius/internal/version.Version=${VERSION} \
        -X github.com/score-spec/score-radius/internal/version.GitCommit=${GIT_COMMIT} \
        -X github.com/score-spec/score-radius/internal/version.BuildDate=${BUILD_DATE}" \
    -o /usr/local/bin/score-radius ./cmd/score-radius

# We can use static since we don't rely on any linux libs or state, but we need ca-certificates to connect to https/oci with the init command.
FROM dhi.io/static:20260413-alpine3.23@sha256:a00ce6c7df24f0e0b98977cc8511cc072960289d0904af36f528cc5ee0189e9c

# Set the current working directory inside the container.
WORKDIR /score-radius

# Copy the binary from the builder image.
COPY --from=builder /usr/local/bin/score-radius /usr/local/bin/score-radius

# Run the binary.
ENTRYPOINT ["/usr/local/bin/score-radius"]