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
package ssh_server

import (
	"fmt"

	"github.com/gliderlabs/ssh"

	"sec-dev-in-action-src/honeypot/server/logger"
	"sec-dev-in-action-src/honeypot/server/pusher"
	"sec-dev-in-action-src/honeypot/server/util"
)

func StartSsh(addr string, flag bool) error {
	ssh.Handle(func(s ssh.Session) {
		// 送佛送到西，通过ssh蜜罐再指引黑客去下一个蜜罐打卡。
		s.Write([]byte(fmt.Sprintf("您的来源IP:%v不在可信列表范围内，"+
			"按公司的安全规范，请先登录跳板机（jumper.sec.lu），再用跳板机登录服务器。\n", s.RemoteAddr())))
	})

	passwordOpt := ssh.PasswordAuth(func(ctx ssh.Context, password string) bool {
		result := false

		if ctx.User() == "root" && password == "123456" {
			result = true
		}

		if flag {
			localAddr := ctx.LocalAddr().String()
			remoteAddr := ctx.RemoteAddr().String()
			rawIp, ProxyAddr, timeStamp := util.GetRawIp(remoteAddr, localAddr)
			logger.Log.Warningf("timestamp: %v, rawIp: %v, proxyAddr: %v, user: %v, password: %v",
				timeStamp, rawIp, ProxyAddr, ctx.User(), password)

			var message pusher.HoneypotMessage
			message.Timestamp = timeStamp
			message.RawIp = rawIp
			message.ProxyAddr = ProxyAddr.String()
			message.User = ctx.User()
			message.Password = password

			strMessage, _ := message.Build()
			logger.Log.Info(strMessage)
			_ = message.Send()
		}

		return result
	})

	logger.Log.Warningf("starting ssh service on %v", addr)
	err := ssh.ListenAndServe(addr, nil, passwordOpt)
	return err
}
