package main

import (
	"os"
	"runtime"

	"github.com/urfave/cli"

	"sec-dev-in-action-src/scanner/tcp-scanner-final/cmd"
)

func main() {
	app := cli.NewApp()
	app.Name = "port_scanner"
	app.Author = "netxfly"
	app.Email = "x@xsec.io"
	app.Version = "2020/3/8"
	app.Usage = "tcp syn/connect port scanner"
	app.Commands = []cli.Command{cmd.Scan}
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	err := app.Run(os.Args)
	_ = err
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
