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

package routers

import (
	"gopkg.in/macaron.v1"

	"sec-dev-in-action-src/proxy-honeypot/server/log"
	"sec-dev-in-action-src/proxy-honeypot/server/models"
	"sec-dev-in-action-src/proxy-honeypot/server/settings"
	"sec-dev-in-action-src/proxy-honeypot/server/util"

	"encoding/json"
	"strings"
)

func Index(ctx *macaron.Context) {
	_ = ctx.Req.ParseForm()
	log.Logger.Info(ctx.Req.Form)
	_, _ = ctx.Write([]byte("test only"))
}

func RecvData(ctx *macaron.Context) {
	_ = ctx.Req.ParseForm()
	timestamp := ctx.Req.Form.Get("timestamp")
	secureKey := ctx.Req.Form.Get("secureKey")
	data := ctx.Req.Form.Get("data")
	agentHost := ctx.Req.Form.Get("hostname")

	headers := ctx.Req.Header

	// get remote ips
	realIp := headers["X-Forwarded-For"]
	ips := make([]string, 0)
	if len(realIp) > 0 {
		t := strings.Split(realIp[0], ",")
		for _, ip := range t {
			sliceIp := strings.Split(ip, ".")
			if len(sliceIp) == 4 {
				ips = append(ips, strings.TrimSpace(ip))
			}

		}
	} else {
		ips = append(ips, ctx.Req.RemoteAddr)
	}

	mySecretKey := util.MakeSign(timestamp, settings.SECRET)
	if secureKey == mySecretKey {
		var h models.HttpRecord
		err := json.Unmarshal([]byte(data), &h)
		// log.Logger.Info(resp, err)
		agentIp := util.Address2Ip(ctx.Req.RemoteAddr)
		if err == nil {
			if len(ips) > 0 {
				agentIp = ips[0]
			}
			record := models.NewRecord(agentIp, agentHost, h)
			err = record.Insert()
			log.Logger.Infof("record: %v, err: %v", record, err)
		}
	} else {
		_, _ = ctx.Write([]byte("error"))
	}
}
