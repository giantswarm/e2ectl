package create

import (
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

const (
	flagName   = "name"
	flagRetain = "retain"
)

type flag struct {
	Name   string
	Retain bool
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.Name, flagName, "kind", `Name of e2e cluster, defaults to "kind".`)
	cmd.Flags().BoolVar(&f.Retain, flagRetain, true, `Retain nodes for debugging when cluster creation fails, defaults to true.`)
}

func (f *flag) Validate() error {
	if f.Name == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagName)
	}

	return nil
}
