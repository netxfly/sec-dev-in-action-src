package main

import (
	"os"
	"runtime"

	"sec-dev-in-action-src/scanner/password_crack/cmd"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "password-crack"
	app.Author = "netxfly"
	app.Email = "x@xsec.io"
	app.Version = "2020/3/11"
	app.Usage = "Weak password scanner"
	app.Commands = []cli.Command{cmd.Scan}
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	err := app.Run(os.Args)
	_ = err
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
