FROM --platform=$BUILDPLATFORM dhi.io/golang:1.27.2-alpine3.24-dev@sha256:9a9160e3c89a37572e35271c65a2b38a2d1671c8ab0d6abee19ce3d6cecb9be9 AS builder

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
FROM dhi.io/static:20260909-alpine3.24@sha256:069a5570611dffc0100e35cd0954ced49fb1299c99f25797d778b9ea73868b9c

# Set the current working directory inside the container.
WORKDIR /score-radius

# Copy the binary from the builder image.
COPY --from=builder /usr/local/bin/score-radius /usr/local/bin/score-radius

# Run the binary.
ENTRYPOINT ["/usr/local/bin/score-radius"]
