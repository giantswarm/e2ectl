[![CircleCI](https://circleci.com/gh/giantswarm/e2ectl.svg?&style=shield)](https://circleci.com/gh/giantswarm/e2ectl)

# e2ectl

Command line tool for managing Kubernetes clusters for use in e2e (integration)
tests. Uses [kind] (Kubernetes in Docker) to do this.

## Design Goals

- `e2ectl` is our CLI and is an opinionated wrapper of [kind].
- Ensures the `kind` clusters we run locally are identical to the clusters we
run in Circle CI.
- We use `e2ectl` in `config.yml` in our Circle CI jobs to provision clusters.

## Installation

This project uses Go modules. Be sure to have it outside your `$GOPATH` or
set `GO111MODULE=on` environment variable. Then regular `go install` should do
the trick. Alternatively the following one-liner may help. 

```sh
GO111MODULE=on go install -ldflags "-X main.gitCommit=$(git rev-parse HEAD)" .
```

## Usage

- Please check `e2ectl -h` for details on all functions.

### Create cluster

- Create a new `kind` cluster.
- Uses our own [retagged-image].

```bash
$ e2ectl cluster create -h

Create cluster for use in e2e tests.

Usage:
  e2ectl cluster create [flags]

Flags:
  -h, --help               help for create
      --image string       Kubernetes image to run. (default "quay.io/giantswarm/kind-node")
      --list-versions      List available Kubernetes versions.
      --name string        Name of e2e cluster. (default "kind")
      --retain             Retain nodes for debugging when cluster creation fails. (default true)
      --version string     Kubernetes version to run. (default "v1.14.2")
      --worker-count int   Number of worker nodes to provision.
```

### Delete cluster

- Delete a `kind` cluster.

```bash
$ e2ectl cluster delete -h

Delete cluster.

Usage:
  e2ectl cluster delete [flags]

Flags:
  -h, --help          help for delete
      --name string   Name of e2e cluster. (default "kind")
```

### Get cluster's kubeconfig file path

- Outputs the kubeconfig file path.

```bash
$ e2ectl kubeconfig path -h

Retrieves cluster kubeconfig file path for use in e2e tests.

Usage:
  e2ectl kubeconfig path [flags]

Flags:
  -h, --help          help for path
      --name string   Name of e2e cluster. (default "kind")

```

## Testing locally

### Use kubectl

```bash
kubectl --kubeconfig=$(kind get kubeconfig-path --name="kind") cluster-info
```

### Connect an operator

- Note: this only works for operators built using our [operatorkit].

```bash
go run main.go daemon \
    --service.kubernetes.incluster=false \
    --service.kubernetes.kubeconfig="$(cat $(kind get kubeconfig-path --name='kind'))"
```

Note: this is for [app-operator]. Other operators may need more flags.

## License

e2ectl is under the Apache 2.0 license. See the [LICENSE](LICENSE) file
for details.

## Credit

- https://github.com/kubernetes-sigs/kind

[app-operator]: https://github.com/giantswarm/app-operator
[kind]: https://kind.sigs.k8s.io/
[operatorkit]: https://github.com/giantswarm/operatorkit
[retagged-image]: https://quay.io/repository/giantswarm/kind-node
