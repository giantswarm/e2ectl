package create

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster"
)

const (
	listVersionsURL = "https://quay.io/api/v1/repository/giantswarm/kind-node/tag/"
	waitForReady    = 2 * time.Minute
)

type runner struct {
	fileSystem afero.Fs
	flag       *flag
	logger     micrologger.Logger
	stdout     io.Writer
	stderr     io.Writer
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

	var provider *cluster.Provider
	{
		provider = cluster.NewProvider()
	}

	{
		// Check if the cluster name already exists.
		n, err := provider.ListNodes(r.flag.Name)
		if err != nil {
			return err
		}
		if len(n) != 0 {
			return microerror.Maskf(invalidFlagError, "cluster with name %#q already exists", r.flag.Name)
		}
	}

	var configData string

	{
		// Define nodes.
		nodes := []KindNode{
			{
				Type: "control-plane",
			},
		}

		for i := 0; i < r.flag.WorkerCount; i++ {
			node := KindNode{
				Type: "worker",
			}

			nodes = append(nodes, node)
		}

		data := KindConfig{
			Nodes: nodes,
		}

		// Render kind config file.
		t, err := template.New("kind").Parse(kindConfigTemplate)
		if err != nil {
			return microerror.Mask(err)
		}

		b := new(bytes.Buffer)
		err = t.Execute(b, data)
		if err != nil {
			return microerror.Mask(err)
		}

		configData = b.String()
	}

	var configFile afero.File

	defer func() {
		err := r.fileSystem.Remove(configFile.Name())
		if err != nil {
			r.logger.LogCtx(ctx, "level", "error", "message", fmt.Sprintf("deletion of %q failed", configFile.Name()), "stack", fmt.Sprintf("%#v", err))
		}
	}()

	{
		configFile, err = afero.TempFile(r.fileSystem, "", "kind-config.yaml")
		if err != nil {
			return microerror.Mask(err)
		}

		_, err = configFile.WriteString(configData)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var image string
	{
		image = fmt.Sprintf("%s:%s", r.flag.Image, r.flag.Version)
	}

	{
		fmt.Printf("creating cluster %#q\n", r.flag.Name)

		err = provider.Create(r.flag.Name,
			cluster.CreateWithConfigFile(configFile.Name()),
			cluster.CreateWithNodeImage(image),
			cluster.CreateWithRetain(r.flag.Retain),
			cluster.CreateWithWaitForReady(waitForReady),
		)
		if err != nil {
			return microerror.Mask(err)
		}

		fmt.Printf("created cluster %#q\n", r.flag.Name)
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
