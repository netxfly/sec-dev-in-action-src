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
	"net/url"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"

	"sec-dev-in-action-src/honeypot/agent/logger"
	"sec-dev-in-action-src/honeypot/agent/settings"
)

var (
	device      string
	snapshotLen int32 = 1024
	promiscuous bool
	err         error
	handle      *pcap.Handle

	filter  = ""
	timeout = time.Duration(3)

	ApiUrl    string
	SecureKey string

	Ips []string

	ApiIp     string
	SensorIps []string
)

func init() {
	device = settings.InterfaceName
	ApiUrl = settings.ManagerUrl
	SecureKey = settings.SecKey

	Ips, err = GetIpList(device)

	urlParsed, err := url.Parse(ApiUrl)
	if err == nil {
		apiHost := urlParsed.Host
		ApiIp, _ = Host2Ip(strings.Split(apiHost, ":")[0])
		SensorIps = Ips
	}

	logger.Log.Infof("local address: %v, apiIp: %v", SensorIps, ApiIp)

	// 给hookHttp添加hook
	hookHttp, err := logger.NewHttpHook()
	if err == nil {
		logger.LogReport.Logger.AddHook(hookHttp)
	}
}

func Start() {
	// Open device
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer handle.Close()
	err = handle.SetBPFFilter(filter)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	ProcessPacket(packetSource.Packets())
}
