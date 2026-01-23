FROM --platform=$BUILDPLATFORM dhi.io/golang:1.25.6-alpine3.22-dev@sha256:1627f7982c2888a3ef84791328a6ee3b8891c16278673627b090740919915b0a AS builder

ARG VERSION

# Set the current working directory inside the container.
WORKDIR /go/src/github.com/score-spec/score-radius

# Copy just the module bits
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project and build it.
COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -ldflags="-s -w -X github.com/score-spec/score-radius/internal/version.Version=${VERSION}" -o /usr/local/bin/score-radius ./cmd/score-radius

# We can use static since we don't rely on any linux libs or state, but we need ca-certificates to connect to https/oci with the init command.
FROM dhi.io/static:20250911-alpine3.22@sha256:6dc61f258412cea484153ee047bd8dcbecc6dc15941befa28a2b82696804b41b

# Set the current working directory inside the container.
WORKDIR /score-radius

# Copy the binary from the builder image.
COPY --from=builder /usr/local/bin/score-radius /usr/local/bin/score-radius

# Run the binary.
ENTRYPOINT ["/usr/local/bin/score-radius"]