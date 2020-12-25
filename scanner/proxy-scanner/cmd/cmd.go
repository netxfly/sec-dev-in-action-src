package cmd

import (
	"github.com/urfave/cli"

	"sec-dev-in-action-src/scanner/proxy-scanner/proxy"
)

var Scan = cli.Command{
	Name:        "scan",
	Usage:       "start to scan proxy",
	Description: "start to scan proxy",
	Action:      proxy.Scan,
	Flags: []cli.Flag{
		boolFlag("debug, d", "debug mode"),
		intFlag("scan_num, n", 100, "scan num"),
		intFlag("timeout, t", 5, "timeout"),
		stringFlag("filename, f", "iplist.txt", "filename"),
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
