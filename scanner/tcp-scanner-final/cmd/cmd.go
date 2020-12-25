package cmd

import (
	"github.com/urfave/cli"

	"sec-dev-in-action-src/scanner/tcp-scanner-final/util"
)

// ./main --iplist ip_list --port port_list --mode syn  --timeout 2 --concurrency 10
var Scan = cli.Command{
	Name:        "scan",
	Usage:       "start to scan port",
	Description: "start to scan port",
	Action:      util.Scan,
	Flags: []cli.Flag{
		stringFlag("iplist, i", "", "ip list"),
		stringFlag("port, p", "", "port list"),
		stringFlag("mode, m", "", "scan mode"),
		intFlag("timeout, t", 3, "timeout"),
		intFlag("concurrency, c", 1000, "concurrency"),
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
