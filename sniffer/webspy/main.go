package main

import (
	"os"
	"runtime"

	"github.com/urfave/cli"

	"sec-dev-in-action-src/sniffer/webspy/cmd"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "webSpy"
	app.Author = "netxfly"
	app.Email = "x@xsec.io"
	app.Version = "2020/5/16"
	app.Usage = "webSpy, Support local and arp spoof mode"
	app.Commands = []cli.Command{cmd.Start}
	app.Flags = append(app.Flags, cmd.Start.Flags...)
	_ = app.Run(os.Args)
}
