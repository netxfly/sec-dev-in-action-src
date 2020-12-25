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

package sniffer

import (
	"net"

	"github.com/google/gopacket/pcap"

	"sec-dev-in-action-src/honeypot/agent/models"
	"sec-dev-in-action-src/honeypot/agent/vars"
)

func GetIpList(deviceName string) (ips []string, err error) {
	devices, err := pcap.FindAllDevs()
	if err == nil {
		for _, device := range devices {
			if device.Name == deviceName {
				for _, addr := range device.Addresses {
					if addr.IP.To4() != nil {
						ips = append(ips, addr.IP.To4().String())
					}
				}
			}
		}
	}

	return ips, err
}

func Host2Ip(host string) (ip string, err error) {
	addr, err := net.LookupHost(host)
	if len(addr) > 0 {
		ip = addr[0]
	}

	return ip, err
}

func SliceContainsString(slice []string, str string) bool {
	m := make(map[string]bool)
	for _, v := range slice {
		m[v] = true
	}
	_, ok := m[str]
	return ok
}

func IsInWhite(conn *models.ConnectionInfo) (result bool) {

	if SliceContainsString(vars.HoneypotPolicy.WhiteIps, conn.SrcIp) ||
		SliceContainsString(vars.HoneypotPolicy.WhiteIps, conn.DstIp) ||
		SliceContainsString(vars.HoneypotPolicy.WhitePorts, conn.SrcPort) ||
		SliceContainsString(vars.HoneypotPolicy.WhitePorts, conn.DstPort) {
		result = true
	}

	return result
}
