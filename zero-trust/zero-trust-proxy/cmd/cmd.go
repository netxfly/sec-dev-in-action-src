package cmd

import (
	"github.com/urfave/cli"

	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/util"
)

var Serve = cli.Command{
	Name:        "serve",
	Usage:       "start sec proxy",
	Description: "start sec proxy",
	Action:      util.Start,
	Flags: []cli.Flag{
		boolFlag("debug, d", "debug mode"),
		stringFlag("config, c", "config", ""),
	},
}

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func boolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

func intFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}
