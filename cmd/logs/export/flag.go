package export

import (
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

const (
	flagLogsDir = "logs-dir"
	flagName    = "name"
)

type flag struct {
	LogsDir string
	Name    string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.LogsDir, flagLogsDir, "", `Logs are exported to this directory.`)
	cmd.Flags().StringVar(&f.Name, flagName, "kind", `Name of e2e cluster.`)
}

func (f *flag) Validate() error {
	if f.LogsDir == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagLogsDir)
	}
	if f.Name == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagName)
	}

	return nil
}
