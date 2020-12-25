/*

Copyright (c) 2017 xsec.io

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THEq
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/

package main

import (
	"os"
	"path/filepath"
	"runtime"

	"sec-dev-in-action-src/proxy-honeypot/agent/cmd"
	"sec-dev-in-action-src/proxy-honeypot/agent/util"
	"sec-dev-in-action-src/proxy-honeypot/agent/vars"

	"github.com/urfave/cli"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	vars.CurrentDir = util.GetCurDir()
	vars.CaKey = filepath.Join(vars.CurrentDir, "./certs/ca.key")
	vars.CaCert = filepath.Join(vars.CurrentDir, "./certs/ca.cert")

	// log.Logger.Infof("dir: %v, caKey: %v, caCert: %v", vars.CurrentDir, vars.CaKey, vars.CaCert)
}

func main() {
	app := cli.NewApp()
	app.Usage = "x-proxy agent"
	app.Version = "0.1"
	app.Author = "netxfly"
	app.Email = "x@xsec.io"
	app.Commands = []cli.Command{cmd.Serve,}
	app.Flags = append(app.Flags, cmd.Serve.Flags...)
	_ = app.Run(os.Args)
}
