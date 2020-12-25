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

	"sec-dev-in-action-src/traffic-analysis/server/models"
)

func Index(ctx *macaron.Context) {
	info, _ := models.ListEvilInfo()
	ctx.Data["info"] = info
	ctx.HTML(200, "index")
}

func HttpReq(ctx *macaron.Context) {
	info, _ := models.ListEvilHttpReq()
	ctx.Data["info"] = info
	ctx.HTML(200, "http_req")
}

func Dns(ctx *macaron.Context) {
	info, _ := models.ListEvilDns()
	ctx.Data["info"] = info
	ctx.HTML(200, "dns")
}
