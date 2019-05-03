package create

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cluster/config"
	"sigs.k8s.io/kind/pkg/cluster/create"
)

const (
	envE2EKubeconfig = "E2E_KUBECONFIG"
	listVersionsURL  = "https://registry.hub.docker.com/v1/repositories/kindest/node/tags"
	nodeImage        = "kindest/node"
	waitForReady     = 2 * time.Minute
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

	if r.flag.ListVersions {
		versions, err := listVersions()
		if err != nil {
			return microerror.Mask(err)
		}

		fmt.Println(versions)

		return nil
	}

	kindCtx := cluster.NewContext(r.flag.Name)
	cfg := &config.Cluster{}

	if r.flag.Version != "" {
		image := fmt.Sprintf("%s:%s", nodeImage, r.flag.Version)

		// Apply image override to all the Nodes defined in Config
		cfg.Nodes = []config.Node{
			{
				Role:  config.ControlPlaneRole,
				Image: image,
			},
		}

		err := cfg.Validate()
		if err != nil {
			return microerror.Mask(err)
		}
	}

	{
		err = kindCtx.Create(cfg, create.Retain(r.flag.Retain), create.WaitForReady(waitForReady))
		if err != nil {
			return microerror.Mask(err)
		}
	}

	{
		kubeconfigPath := kindCtx.KubeConfigPath()
		err := os.Setenv(envE2EKubeconfig, kubeconfigPath)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func listVersions() ([]string, error) {
	resp, err := http.Get(listVersionsURL)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	defer resp.Body.Close()

	data := []struct {
		Name string `json:"name"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var versions []string
	for _, version := range data {
		versions = append(versions, version.Name)
	}

	return versions, nil
}
