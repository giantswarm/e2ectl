package cluster

import (
	"io"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"

	"github.com/giantswarm/e2ectl/cmd/cluster/create"
	"github.com/giantswarm/e2ectl/cmd/cluster/delete"
)

const (
	name        = "cluster"
	description = "Manage clusters."
)

type Config struct {
	Logger micrologger.Logger
	Stderr io.Writer
	Stdout io.Writer
}

func New(config Config) (*cobra.Command, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	var err error

	var createCmd *cobra.Command
	{
		c := create.Config{
			Logger: config.Logger,
			Stderr: config.Stderr,
			Stdout: config.Stdout,
		}

		createCmd, err = create.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var deleteCmd *cobra.Command
	{
		c := delete.Config{
			Logger: config.Logger,
			Stderr: config.Stderr,
			Stdout: config.Stdout,
		}

		deleteCmd, err = delete.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	f := &flag{}

	r := &runner{
		flag:   f,
		logger: config.Logger,
		stderr: config.Stderr,
		stdout: config.Stdout,
	}

	c := &cobra.Command{
		Use:   name,
		Short: description,
		Long:  description,
		RunE:  r.Run,
	}

	f.Init(c)

	c.AddCommand(createCmd)
	c.AddCommand(deleteCmd)

	return c, nil
}
