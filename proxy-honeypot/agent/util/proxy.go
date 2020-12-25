/*

Copyright (c) 2018 sec.lu

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

package util

import (
	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/html"
	"github.com/urfave/cli"

	"sec-dev-in-action-src/proxy-honeypot/agent/log"
	"sec-dev-in-action-src/proxy-honeypot/agent/modules"
	"sec-dev-in-action-src/proxy-honeypot/agent/settings"
	"sec-dev-in-action-src/proxy-honeypot/agent/vars"

	"fmt"
	"net/http"
)

func init() {
	cfg := settings.Cfg
	sec := cfg.Section("proxy")
	vars.ProxyHost = sec.Key("HOST").MustString("")
	vars.ProxyPort = sec.Key("PORT").MustInt(1080)
	vars.DebugMode = sec.Key("DEBUG").MustBool(false)

}

func Start(ctx *cli.Context) {
	if ctx.IsSet("debug") {
		vars.DebugMode = ctx.Bool("debug")
	}

	if ctx.IsSet("port") {
		vars.ProxyPort = ctx.Int("port")
	}

	err := SetCA()
	log.Logger.Infof("caKey: %v, caCert: %v, set ca err: %v", vars.CaKey, vars.CaCert, err)

	proxy := goproxy.NewProxyHttpServer()
	log.Logger.Infof("proxy Start success, Listening on %v:%v ", vars.ProxyHost, vars.ProxyPort)

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest().DoFunc(modules.ReqHandlerFunc)
	proxy.OnResponse(goproxy_html.IsWebRelatedText).DoFunc(modules.RespHandlerFunc)

	proxy.Verbose = vars.DebugMode

	log.Logger.Info(http.ListenAndServe(fmt.Sprintf("%v:%v", vars.ProxyHost, vars.ProxyPort), proxy))
}
