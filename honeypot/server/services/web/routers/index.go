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

package routers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"sec-dev-in-action-src/honeypot/server/logger"
	"sec-dev-in-action-src/honeypot/server/pusher"
	"sec-dev-in-action-src/honeypot/server/util"
)

var httpAddr = "127.0.0.1:8000"

func IndexHandle(ctx *gin.Context) {
	_, ok := ctx.Get("flag")
	_ = ctx.Request.ParseForm()
	params := ctx.Request.Form

	remoteAddr := ctx.Request.RemoteAddr
	host := ctx.Request.Host

	body := make([]byte, 0)
	n, err := ctx.Request.Body.Read(body)
	logger.Log.Infof("n: %v, err: %v", n, err)
	if ok {
		rawIp, ProxyAddr, timeStamp := util.GetRawIp(remoteAddr, httpAddr)
		logger.Log.Warnf("rawIp: %v, proxyAddr: %v, timestamp: %v", rawIp, ProxyAddr, timeStamp)
		var message pusher.HoneypotMessage
		message.Timestamp = timeStamp
		message.RawIp = rawIp
		message.ProxyAddr = ProxyAddr.String()

		data := make(map[string]interface{})
		data["body"] = body

		message.Data = data
		strMessage, _ := message.Build()
		logger.Log.Info(strMessage)
		_ = message.Send()
	}

	ctx.String(http.StatusOK, fmt.Sprintf("Hello, World! \nremote_addr: %v, host: %v, param: %v, body: %v\n",
		remoteAddr, host, params, string(body)))
}
