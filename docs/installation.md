# Installation of `score-radius`

You can install the `score-radius` CLI in a variety of ways:
- [Homebrew](#homebrew)
- [Go](#go)
- [Docker](#docker)
- [Manual download](#manual-download)

## Homebrew

_Prerequisites: You must have [brew](https://brew.sh/) installed._

```bash
brew install score-spec/tap/score-radius
```

## Go

_Prerequisites: You must have [Go](https://go.dev/dl/) installed._

```bash
go install -v github.com/score-spec/score-radius/cmd/score-radius@latest
```

## Docker

_Prerequisites: You must have [Docker](https://docs.docker.com/get-docker/) installed._

```bash
docker run --rm -it ghcr.io/score-spec/score-radius:latest
```

If you want to run `score-radius` with the `--help` flag to view the available options, you would run the following command.
```bash
docker run --rm -it ghcr.io/score-spec/score-radius:latest --help
```

If you want to run `score-radius` with the `init` subcommand to initialize your local working directory, you would run the following command.
```bash
docker run --rm -it -v .:/score-radius ghcr.io/score-spec/score-radius:latest init
```

If you want to run `score-radius` as an unprivileged container, you would run the following command.

```bash
docker run --rm -it -v .:/score-radius --read-only --cap-drop=ALL --user=65532 ghcr.io/score-spec/score-radius:latest init
```

## Manual download

You can see the list of available files per OS and architecture in the [`score-radius`'s releases page](https://github.com/score-spec/score-radius/releases).


Here is an example showing how to install the latest version of `score-radius` via `wget` on Linux, you can adapt this based on your own environment:
```bash
VERSION=$(curl -sL https://api.github.com/repos/score-spec/score-radius/releases/latest | jq -r .tag_name)
OS=linux
ARCH=amd64

wget https://github.com/score-spec/score-radius/releases/download/${VERSION}/score-radius_${VERSION:1}_${OS}_${ARCH}.tar.gz

tar -xvf score-radius_${VERSION:1}_${OS}_${ARCH}.tar.gz

chmod +x score-radius
sudo mv score-radius /usr/local/bin
```
