package main

import (
	"context"
	"fmt"
	"os"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/giantswarm/e2ectl/cmd"
)

var (
	gitCommit = "n/a"
	source    = "https://github.com/giantswarm/e2ectl"
)

func main() {
	err := mainWithError()
	if err != nil {
		panic(fmt.Sprintf("%#v\n", err))
	}
}

func mainWithError() error {
	var err error
	ctx := context.Background()

	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	fileSystem := afero.NewOsFs()

	var rootCommand *cobra.Command
	{
		c := cmd.Config{
			FileSystem: fileSystem,
			Logger:     logger,

			GitCommit: gitCommit,
			Source:    source,
		}

		rootCommand, err = cmd.New(c)
	}

	err = rootCommand.Execute()
	if err != nil {
		logger.LogCtx(ctx, "level", "error", "message", "failed to execute command", "stack", fmt.Sprintf("%#v", err))
		os.Exit(1)
	}

	return nil
}
