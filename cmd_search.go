package main

import (
	"errors"

	"github.com/urfave/cli"
)

func cmdSearch(c *cli.Context) error {
	pluginName := c.Args().First()
	if pluginName == "" {
		return errors.New("plugin name required")
	}

	if err := searchPlugin(pluginName); err != nil {
		fatal("Error: Fail search plugin.")
	}

	return nil
}
