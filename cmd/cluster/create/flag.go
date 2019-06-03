package create

import (
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"

	"sigs.k8s.io/kind/pkg/cluster/config/defaults"
)

const (
	flagListVersions = "list-versions"
	flagName         = "name"
	flagRetain       = "retain"
	flagVersion      = "version"
	flagWorkerCount  = "worker-count"
)

type flag struct {
	ListVersions bool
	Name         string
	Retain       bool
	Version      string
	WorkerCount  int
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&f.ListVersions, flagListVersions, false, `List available Kubernetes version.`)
	cmd.Flags().StringVar(&f.Name, flagName, "kind", `Name of e2e cluster.`)
	cmd.Flags().BoolVar(&f.Retain, flagRetain, true, `Retain nodes for debugging when cluster creation fails.`)
	cmd.Flags().IntVar(&f.WorkerCount, flagWorkerCount, 0, `Number of worker nodes to provision.`)

	var defaultVersion string
	s := strings.Split(defaults.Image, ":")
	if len(s) >= 2 {
		defaultVersion = s[1]
	}
	cmd.Flags().StringVar(&f.Version, flagVersion, defaultVersion, `Kubernetes version to run.`)
}

func (f *flag) Validate() error {
	if f.Name == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagName)
	}

	return nil
}
