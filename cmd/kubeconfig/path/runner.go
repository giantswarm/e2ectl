package path

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

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

	var provider *cluster.Provider
	{
		provider = cluster.NewProvider()
	}

	{
		// Check if the cluster name exists.
		n, err := provider.ListNodes(r.flag.Name)
		if err != nil {
			return err
		}
		if len(n) == 0 {
			return microerror.Maskf(invalidFlagError, "cluster %#q does not exist", r.flag.Name)
		}
	}

	var kubeconfigPath string
	{
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		kubeconfigPath = path.Join(dir, r.flag.Kubeconfig)
	}

	{
		err = provider.ExportKubeConfig(r.flag.Name, kubeconfigPath)
		if err != nil {
			return err
		}
	}

	fmt.Println(kubeconfigPath)

	return nil
}
