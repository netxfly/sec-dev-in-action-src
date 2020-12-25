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

import "time"

type (
	ConnectionInfo struct {
		Protocol string `json:"protocol,omitempty"`
		SrcIp    string `json:"src_ip,omitempty"`
		SrcPort  string `json:"src_port,omitempty"`
		DstIp    string `json:"dst_ip,omitempty"`
		DstPort  string `json:"dst_port,omitempty"`
		IsHttp   bool   `json:"is_http,omitempty"`
	}

	PacketInfo struct {
		ConnInfo *ConnectionInfo `json:"conn_info"`
		Time     time.Time       `json:"time"`
	}
)

func (p *PacketInfo) Insert() error {
	err := CollectionPacket.Insert(p)
	return err
}
