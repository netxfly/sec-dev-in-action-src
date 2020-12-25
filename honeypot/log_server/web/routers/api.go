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
	"encoding/json"
	"fmt"
	"log"

	"sec-dev-in-action-src/honeypot/log_server/models"
	"sec-dev-in-action-src/honeypot/log_server/settings"
	"sec-dev-in-action-src/honeypot/log_server/util"

	"github.com/gin-gonic/gin"
)

func PacketHandle(ctx *gin.Context) {
	timestamp := ctx.PostForm("timestamp")
	secureKey := ctx.PostForm("secureKey")
	data := ctx.PostForm("data")
	remoteAddr := ctx.Request.RemoteAddr
	mySecureKey := util.MD5(fmt.Sprintf("%v%v", timestamp, settings.KEY))
	packetInfo := models.PacketInfo{}
	log.Printf("data: %v\n", data)
	if mySecureKey == secureKey {
		err := json.Unmarshal([]byte(data), &packetInfo)
		if err == nil {
			err := packetInfo.Insert()
			fmt.Printf("remoteAddr: %v, packetInfo: %v, err: %v\n", remoteAddr, packetInfo, err)
		}

		ctx.JSON(200, "ok")
	} else {
		ctx.JSON(200, "err")
	}
}

func ServiceHandle(ctx *gin.Context) {
	timestamp := ctx.PostForm("timestamp")
	secureKey := ctx.PostForm("secureKey")
	data := ctx.PostForm("data")
	remoteAddr := ctx.Request.RemoteAddr
	mySecureKey := util.MD5(fmt.Sprintf("%v%v", timestamp, settings.KEY))

	log.Printf("data: %v\n", data)
	if secureKey == mySecureKey {
		var message models.HoneypotMessage
		err := json.Unmarshal([]byte(data), &message)
		if err == nil {
			err := message.Insert()
			fmt.Printf("remoteAddr: %v, message: %v, err: %v\n", remoteAddr, message, err)
		}
	}
}
