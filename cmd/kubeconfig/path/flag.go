package path

import (
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster"
)

const (
	flagKubeconfig = "kubeconfig"
	flagName       = "name"
)

type flag struct {
	Kubeconfig string
	Name       string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.Kubeconfig, flagKubeconfig, "kind-kubeconfig", `Name of kubeconfig file.`)
	cmd.Flags().StringVar(&f.Name, flagName, cluster.DefaultName, `Name of e2e cluster.`)
}

func (f *flag) Validate() error {
	if f.Name == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagName)
	}

	return nil
}
