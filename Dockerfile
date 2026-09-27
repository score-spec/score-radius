FROM --platform=$BUILDPLATFORM dhi.io/golang:1.27.1-alpine3.24-dev@sha256:0fbbb101cb3c451453aa0d3e7a87c378c6bd784d8cfbfd4958f97dd3bd0197e6 AS builder

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
FROM dhi.io/static:20260909-alpine3.24@sha256:296ab7284ac616e1f03b9ae929852b968315242311da974c57de342894276418

# Set the current working directory inside the container.
WORKDIR /score-radius

# Copy the binary from the builder image.
COPY --from=builder /usr/local/bin/score-radius /usr/local/bin/score-radius

# Run the binary.
ENTRYPOINT ["/usr/local/bin/score-radius"]
