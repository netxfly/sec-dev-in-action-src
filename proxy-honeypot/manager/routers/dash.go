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
	"sec-dev-in-action-src/proxy-honeypot/manager/models"

	"gopkg.in/macaron.v1"

	"github.com/go-macaron/session"
)

func Dash(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		totalRecord, totalPassword, err := models.DashTotal()
		passwords, err := models.DashPassword()
		urls, err := models.DashUrls()
		evilIps, err := models.DashIps()
		_ = err

		ctx.Data["total_record"] = totalRecord
		ctx.Data["total_password"] = totalPassword

		ctx.Data["passwords"] = passwords
		ctx.Data["urls"] = urls
		ctx.Data["evil_ips"] = evilIps

		ctx.HTML(200, "dash")
	} else {
		ctx.Redirect("/admin/login/")
	}
}
