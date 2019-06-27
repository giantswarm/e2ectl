package path

import (
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster"
)

const (
	flagName = "name"
)

type flag struct {
	Name string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.Name, flagName, cluster.DefaultName, `Name of e2e cluster.`)
}

func (f *flag) Validate() error {
	if f.Name == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagName)
	}

	return nil
}
