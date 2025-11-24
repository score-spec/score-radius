# score-radius

`score-radius` is a Score implementation of the [Score specification](https://score.dev/) for [Radius](https://radapp.io/).

![](./docs/images/score-radius.png)

With the CLI:
```bash
make build

score-radius --version

score-radius init

score-radius generate score.yaml \
    -a my-app-123 \
    -e my-env-123 \
    -o app.bicep
```

With the container image:
```bash
make build-container

docker run --rm -it -v .:/score-radius score-radius:local init

docker run --rm -it -v .:/score-radius score-radius:local generate score.yaml \
    -a my-app-123 \
    -e my-env-123 \
    -o app.bicep
```


## Docs

- [CLI](./docs/cli.md)
- [Known limitations](./docs/known-limitations.md)