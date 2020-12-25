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

package models

import (
	"encoding/json"
	"time"
)

type (
	ConnectionInfo struct {
		Protocol string `json:"protocol"`
		SrcIp    string `json:"src_ip"`
		SrcPort  string `json:"src_port"`
		DstIp    string `json:"dst_ip"`
		DstPort  string `json:"dst_port"`
		IsHttp   bool   `json:"is_http"`
	}

	PacketInfo struct {
		ConnInfo *ConnectionInfo `json:"conn_info"`
		Time     time.Time       `json:"time"`
	}
)

func NewConnectionInfo(proto string, srcIp string, srcPort string, dstIp string, dstPort string, isHttp bool) (connInfo *ConnectionInfo) {
	return &ConnectionInfo{Protocol: proto, SrcIp: srcIp, SrcPort: srcPort, DstIp: dstIp, DstPort: dstPort, IsHttp: isHttp}
}

func (c *ConnectionInfo) String() (string, error) {
	js, err := json.Marshal(c)
	return string(js), err
}

func NewPacketInfo(info *ConnectionInfo, now time.Time) (ret *PacketInfo) {
	ret = &PacketInfo{ConnInfo: info, Time: now}
	return ret
}

func (p *PacketInfo) String() (string, error) {
	js, err := json.Marshal(p)
	return string(js), err
}
