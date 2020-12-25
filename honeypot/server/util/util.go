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
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"sec-dev-in-action-src/honeypot/server/vars"
)

func GetRawIpByConn(conn net.Conn) (string, net.TCPAddr, int64) {
	remoteAddr := conn.RemoteAddr().String()
	localAddr := conn.LocalAddr().String()
	return GetRawIp(remoteAddr, localAddr)
}

func GetRawIp(remoteAddr, localAddr string) (string, net.TCPAddr, int64) {
	var (
		rawIp     string
		ProxyAddr net.TCPAddr
		timeStamp int64
	)

	k := fmt.Sprintf("%v_%v", remoteAddr, localAddr)
	v, ok := vars.RawIps.Load(k)
	fmt.Printf("k: %v, v: %v, ok: %v\n", k, v, ok)
	if ok {
		value, ok := v.(string)
		if ok {
			t := strings.Split(value, "_")
			if len(t) == 2 {
				rawIp = t[1]
				ProxyAddrStr := t[0]
				tt := strings.Split(ProxyAddrStr, "@")
				if len(tt) == 2 {

					timeStamp, _ = strconv.ParseInt(tt[0], 10, 64)
					proxyIpPort := tt[1]
					ttt := strings.Split(proxyIpPort, ":")
					// fmt.Printf("ttt: %v, len(ttt): %v\n", ttt, len(ttt))
					if len(ttt) == 2 {
						ProxyAddr.IP = StrToIp(ttt[0])
						port, _ := strconv.Atoi(ttt[1])
						ProxyAddr.Port = port
					}
				}
			}
		}
	}

	return rawIp, ProxyAddr, timeStamp
}

func StrToIp(ip string) net.IP {
	return net.ParseIP(ip)
}

func DelExpireIps(timeoutSec int64) {
	vars.RawIps.Range(func(key, value interface{}) bool {
		v, ok := value.(string)
		if ok {
			timestamp := getTimestamp(v)
			if time.Now().Unix()-timestamp >= timeoutSec {
				vars.RawIps.Delete(key)
			}
		}
		return ok
	})
}

func getTimestamp(v string) int64 {
	var timestamp int64
	t := strings.Split(v, "_")
	if len(t) == 2 {
		ProxyAddrStr := t[0]
		tt := strings.Split(ProxyAddrStr, "@")
		if len(tt) == 2 {
			timestamp, _ = strconv.ParseInt(tt[0], 10, 64)
		}
	}
	return timestamp
}
