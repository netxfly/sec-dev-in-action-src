package cmd

import (
	"github.com/urfave/cli"

	"sec-dev-in-action-src/sniffer/webspy/modules"
)

var Start = cli.Command{
	Name:        "start",
	Usage:       "sniff local server",
	Description: "startup sniff on local server",
	Action:      modules.Start,
	Flags: []cli.Flag{
		stringFlag("mode,m", "local", "webspy running mode, local or arp"),
		stringFlag("device,i", "eth0", "device name"),
		stringFlag("host,H", "127.0.0.1", "web server listen address"),
		intFlag("port,p", 4000, "web server listen address"),
		boolFlag("debug, d", "debug mode"),
		stringFlag("target, t", "", "target ip address"),
		stringFlag("gateway, g", "", "gateway ip address"),
		stringFlag("filter,f", "", "setting filters"),
		intFlag("length,l", 1024, "setting snapshot Length"),
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
