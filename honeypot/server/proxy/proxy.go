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

package proxy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"time"

	"sec-dev-in-action-src/honeypot/server/logger"
	"sec-dev-in-action-src/honeypot/server/util"
	"sec-dev-in-action-src/honeypot/server/vars"
)

var (
	sshLocalAddr   = fmt.Sprintf("%v:8022", vars.Config.Proxy.Addr)
	sshBackendAddr = "127.0.0.1:2222"

	mysqlLocalAddr   = fmt.Sprintf("%v:3306", vars.Config.Proxy.Addr)
	mysqlBackendAddr = "127.0.0.1:3366"

	redisLocalAddr   = fmt.Sprintf("%v:6379", vars.Config.Proxy.Addr)
	redisBackendAddr = "127.0.0.1:6380"

	webLocalAddr   = fmt.Sprintf("%v:8080", vars.Config.Proxy.Addr)
	webBackendAddr = "127.0.0.1:8000"
)

// 再加个代理获取下ssh的真实源IP
func StartProxy() {
	// ssh的代理
	go serveProxy(sshLocalAddr, sshBackendAddr)
	// redis的代理
	go serveProxy(redisLocalAddr, redisBackendAddr)
	// web的代理
	go serveProxy(webLocalAddr, webBackendAddr)
	// mysql的代理
	go serveProxy(mysqlLocalAddr, mysqlBackendAddr)

	for {
		time.Sleep(10 * time.Second)
		util.DelExpireIps(300)
	}
}

func serveProxy(localAddr, backendAddr string) {
	lis, err := net.Listen("tcp", localAddr)
	if err != nil {
		return
	}
	defer lis.Close()
	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}

		host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
		_ = port
		remoteIp := net.ParseIP(host)
		if remoteIp == nil {
			continue
		}
		go handleConn(conn, backendAddr)
	}
}

func handlePipe(srcConn net.Conn, dstConn net.Conn, flag bool) error {
	data := make([]byte, 10240)
	for {
		n, err := srcConn.Read(data)
		// logger.Log.Infof("n: %v, err: %v", n, err)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		b := data[:4]
		if flag {
			intIp, err := BytesToInt(b)
			// logger.Log.Infof("intIp: %v, err: %v", intIp, err)
			if err != nil {
				return err
			}
			srcIp := Uint2IP(intIp).String()
			key := fmt.Sprintf("%v_%v", dstConn.LocalAddr(), dstConn.RemoteAddr())
			value := fmt.Sprintf("%v@%v_%v", time.Now().Unix(), srcConn.RemoteAddr(), srcIp)
			vars.RawIps.Store(key, value)
			logger.Log.Infof("srcIp: %v, set key:%v -> value: %v", srcIp, key, value)
			b = data[4:n]
		} else {
			b = data[:n]
		}
		_, _ = dstConn.Write(b)
	}

	return nil
}

func handleConn(srcConn net.Conn, backendAddr string) {
	defer srcConn.Close()
	ip := backendAddr
	dstConn, err := net.DialTimeout("tcp", ip, 5*time.Second)
	if err != nil {
		return
	}
	defer dstConn.Close()

	go func(SrcConn net.Conn, dstConn net.Conn) {
		err := handlePipe(srcConn, dstConn, true)
		_ = err
	}(srcConn, dstConn)

	go func(srcConn net.Conn, dstConn net.Conn) {
		err := handlePipe(dstConn, srcConn, false)
		_ = err
	}(srcConn, dstConn)

	exit := make(chan bool, 1)
	<-exit
}

// IP2Uint 把net.IP转为数值
func IP2Uint(ip net.IP) uint32 {
	b := ip.To4()
	if b == nil {
		return 0
	}

	return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}

// Uint2IP 把数值转为net.IP
func Uint2IP(i uint32) net.IP {
	if i > math.MaxUint32 {
		return nil
	}

	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)

	return ip
}

// 字节转换成整形
func BytesToInt(b []byte) (uint32, error) {
	var x uint32
	bytesBuffer := bytes.NewBuffer(b)
	err := binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x, err
}
