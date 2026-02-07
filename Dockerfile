FROM --platform=$BUILDPLATFORM golang:1.25.7-alpine@sha256:f6751d823c26342f9506c03797d2527668d095b0a15f1862cddb4d927a7a4ced AS builder

ARG VERSION

# Set the current working directory inside the container.
WORKDIR /go/src/github.com/score-spec/score-radius

# Copy just the module bits
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project and build it.
COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -ldflags="-s -w -X github.com/score-spec/score-radius/internal/version.Version=${VERSION}" -o /usr/local/bin/score-radius ./cmd/score-radius

# We can use gcr.io/distroless/static since we don't rely on any linux libs or state, but we need ca-certificates to connect to https/oci with the init command.
FROM gcr.io/distroless/static-debian13:nonroot@sha256:f9f84bd968430d7d35e8e6d55c40efb0b980829ec42920a49e60e65eac0d83fc

# Set the current working directory inside the container.
WORKDIR /score-radius

# Copy the binary from the builder image.
COPY --from=builder /usr/local/bin/score-radius /usr/local/bin/score-radius

# Run the binary.
ENTRYPOINT ["/usr/local/bin/score-radius"]