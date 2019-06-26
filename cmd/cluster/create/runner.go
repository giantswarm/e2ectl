package create

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cluster/config"
	"sigs.k8s.io/kind/pkg/cluster/create"
)

const (
	listVersionsURL = "https://quay.io/api/v1/repository/giantswarm/kind-node/tag/"
	waitForReady    = 2 * time.Minute
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

	var known bool
	{
		known, err = cluster.IsKnown(r.flag.Name)
		if err != nil {
			return err
		}
		if known {
			return microerror.Maskf(invalidFlagError, "cluster with name %#q already exists", r.flag.Name)
		}
	}

	kindCtx := cluster.NewContext(r.flag.Name)

	cfg := &config.Cluster{}
	config.SetDefaults_Cluster(cfg)
	nodes := []config.Node{}

	if r.flag.Version != "" {
		image := fmt.Sprintf("%s:%s", r.flag.Image, r.flag.Version)

		controlPlane := config.Node{
			Image: image,
			Role:  config.ControlPlaneRole,
		}

		nodes = append(nodes, controlPlane)

		for i := 0; i < r.flag.WorkerCount; i++ {
			worker := config.Node{
				Image: image,
				Role:  config.WorkerRole,
			}

			nodes = append(nodes, worker)
		}

		cfg.Nodes = nodes

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

	return nil
}

func listVersions() ([]string, error) {
	resp, err := http.Get(listVersionsURL)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	defer resp.Body.Close()

	data := struct {
		Tags []struct {
			Name string `json:"name"`
		} `json:"tags"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var versions []string
	for _, tag := range data.Tags {
		versions = append(versions, tag.Name)
	}

	return versions, nil
}
