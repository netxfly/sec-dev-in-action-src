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
	"sec-dev-in-action-src/honeypot/agent/logger"
	"sec-dev-in-action-src/honeypot/agent/vars"
	"strings"
	"time"
)

type (
	Pipe struct {
		source net.Conn
		dest   net.Conn
	}

	ForwardPolicy struct {
		LocalPort  int
		TargetHost string
		TargetPort int
	}
)

func processConnection(srcConn net.Conn, targetAddr string) {
	destConn, err := net.DialTimeout("tcp", targetAddr, 3*time.Second)
	if err != nil {
		logger.Log.Errorf("Unable to connect to %s, %v\n", targetAddr, err)
		_ = srcConn.Close()
		return
	}

	// 发送数据
	go func(srcConn net.Conn, destConn net.Conn) {
		err := passThrough(srcConn, destConn, true)
		_ = err
	}(srcConn, destConn)

	// 接收数据
	go func(srcConn net.Conn, destConn net.Conn) {
		err := passThrough(destConn, srcConn, false)
		_ = err
	}(srcConn, destConn)
}

func passThrough(srcConn net.Conn, destConn net.Conn, proxyFlag bool) error {
	remoteIp := GetRemoteIP(srcConn)
	BytesIntIp := IntToBytes(IP2Uint(net.ParseIP(remoteIp)))

	data := make([]byte, 10240)
	for {
		n, err := srcConn.Read(data)
		if err != nil {
			return err
		}
		if err == io.EOF {
			break
		}
		buffer := data[:n]
		if proxyFlag {
			// logger.Log.Infof("remoteIp: %v, bytesIntIp: %v", remoteIp, BytesIntIp)
			buffer = append(BytesIntIp, buffer...)
		}

		_, _ = destConn.Write(buffer)
	}

	return nil
}

func Proxy() {
	forwardPolicy := make([]ForwardPolicy, 0)

	for _, service := range vars.PolicyData.Service {
		forwardPolicy = append(forwardPolicy, ForwardPolicy{LocalPort: service.LocalPort,
			TargetHost: service.BackendHost,
			TargetPort: service.BackendPort},
		)
	}

	for _, item := range forwardPolicy {
		go func(item ForwardPolicy) {
			target := fmt.Sprintf("%v:%v", item.TargetHost, item.TargetPort)
			listener, err := net.Listen("tcp", fmt.Sprintf(":%v", item.LocalPort))
			if err != nil {
				return
			}

			logger.Log.Infof("Forward :%v -> %v:%v", item.LocalPort, item.TargetHost, item.TargetPort)

			for {
				conn, err := listener.Accept()
				if err != nil {
					logger.Log.Errorf("Accept failed, %v\n", err)
					break
				}
				logger.Log.Infof("%v -> %v -> %v", conn.LocalAddr(), conn.RemoteAddr(), target)
				go processConnection(conn, target)
			}
		}(item)
	}

	for {
		time.Sleep(30 * time.Second)
	}

}

func GetRemoteIP(c net.Conn) string {
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.RemoteAddr().String())); err == nil {
		return ip
	}

	return ""
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

// 整形转换为字节型
func IntToBytes(n uint32) []byte {
	x := uint32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}
