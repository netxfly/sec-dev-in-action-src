package main

import (
	"os"
	"runtime"

	"github.com/urfave/cli"

	"sec-dev-in-action-src/scanner/proxy-scanner/cmd"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "proxy scanner"
	app.Author = "netxfly"
	app.Email = "x@xsec.io"
	app.Version = "2020/5/15"
	app.Usage = "A SOCKS4/SOCKS4a/SOCKS5/HTTP/HTTPS proxy scanner"
	app.Commands = []cli.Command{cmd.Scan}
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	_ = app.Run(os.Args)
}
