[![CircleCI](https://circleci.com/gh/giantswarm/e2ectl.svg?&style=shield)](https://circleci.com/gh/giantswarm/e2ectl)

# e2ectl

Command line tool for managing Kubernetes clusters for use in e2e (integration)
tests. Uses [kind] (Kubernetes in Docker) to do this.

## Installation

This project uses Go modules. Be sure to have it outside your `$GOPATH` or
set `GO111MODULE=on` environment variable. Then regular `go install` should do
the trick. Alternatively the following one-liner may help. 

```sh
GO111MODULE=on go install -ldflags "-X main.gitCommit=$(git rev-parse HEAD)" .
```

## Usage

Please check `e2ectl -h` for details on all functions.

## License

e2ectl is under the Apache 2.0 license. See the [LICENSE](LICENSE) file
for details.

## Credit

- https://github.com/kubernetes-sigs/kind

[kind]: https://kind.sigs.k8s.io/
