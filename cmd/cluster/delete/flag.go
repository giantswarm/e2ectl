package delete

import (
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

const (
	flagName = "name"
)

type flag struct {
	Name string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.Name, flagName, "kind", `Name of e2e cluster, defaults to "kind".`)
}

func (f *flag) Validate() error {
	if f.Name == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagName)
	}

	return nil
}
