package delete

import (
	"context"
	"io"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster"
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
		known, err := cluster.IsKnown(r.flag.Name)
		if err != nil {
			return err
		}
		if !known {
			return microerror.Maskf(invalidFlagError, "no cluster with name %#q", r.flag.Name)
		}
	}

	{
		kindCtx := cluster.NewContext(r.flag.Name)

		err = kindCtx.Delete()
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
