package main

import (
	"os"
	"runtime"

	"github.com/urfave/cli"

	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/cmd"
)

func main() {
	app := cli.NewApp()
	app.Usage = "zero-trust-proxy-demo"
	app.Version = "0.1"
	app.Author = "netxfly"
	app.Email = "x@xsec.io"
	app.Commands = []cli.Command{cmd.Serve}
	app.Flags = append(app.Flags, cmd.Serve.Flags...)
	_ = app.Run(os.Args)
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
