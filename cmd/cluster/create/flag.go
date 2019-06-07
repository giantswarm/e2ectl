package create

import (
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

const (
	flagImage        = "image"
	flagListVersions = "list-versions"
	flagName         = "name"
	flagRetain       = "retain"
	flagVersion      = "version"
	flagWorkerCount  = "worker-count"
)

type flag struct {
	Image        string
	ListVersions bool
	Name         string
	Retain       bool
	Version      string
	WorkerCount  int
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.Image, flagImage, "quay.io/giantswarm/kind-node", `Kubernetes image to run.`)
	cmd.Flags().BoolVar(&f.ListVersions, flagListVersions, false, `List available Kubernetes versions.`)
	cmd.Flags().StringVar(&f.Name, flagName, "kind", `Name of e2e cluster.`)
	cmd.Flags().BoolVar(&f.Retain, flagRetain, true, `Retain nodes for debugging when cluster creation fails.`)
	cmd.Flags().IntVar(&f.WorkerCount, flagWorkerCount, 0, `Number of worker nodes to provision.`)

	versions, err := listVersions()
	if err != nil {
		panic("Failed to lookup available versions")
	}
	defaultVersion := versions[0]
	cmd.Flags().StringVar(&f.Version, flagVersion, defaultVersion, `Kubernetes version to run.`)
}

func (f *flag) Validate() error {
	if f.Name == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagName)
	}

	return nil
}
