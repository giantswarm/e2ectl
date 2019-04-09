package create

import (
	"context"
	"io"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cluster/config"
	"sigs.k8s.io/kind/pkg/cluster/create"
)

const (
	waitForReady = 2 * time.Minute
)

type runner struct {
	flag   *flag
	logger micrologger.Logger
	stdout io.Writer
	stderr io.Writer
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := r.flag.Validate()
	if err != nil {
		return microerror.Mask(err)
	}

	err = r.run(ctx, cmd, args)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	var err error

	{
		kindCtx := cluster.NewContext(r.flag.Name)
		cfg := &config.Cluster{}

		err = kindCtx.Create(cfg, create.Retain(r.flag.Retain), create.WaitForReady(waitForReady))
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
