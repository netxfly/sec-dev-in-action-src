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

package sensor

import (
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"sec-dev-in-action-src/traffic-analysis/sensor/misc"
	"sec-dev-in-action-src/traffic-analysis/sensor/settings"
)

var (
	device      string
	snapshotLen int32 = 1024
	promiscuous bool  = true
	err         error
	timeout     time.Duration = pcap.BlockForever
	handle      *pcap.Handle

	DebugMode bool
	filter    = ""

	ApiUrl    string
	SecureKey string
)

func init() {
	device = settings.DeviceName
	DebugMode = settings.DebugMode
	filter = settings.FilterRule

	cfg := settings.Cfg
	sec := cfg.Section("server")
	ApiUrl = sec.Key("API_URL").MustString("")
	SecureKey = sec.Key("API_KEY").MustString("")

}

func Start(ctx *cli.Context) {
	if ctx.IsSet("debug") {
		DebugMode = ctx.Bool("debug")
	}
	if DebugMode {
		misc.Log.Logger.Level = logrus.DebugLevel
	}

	if ctx.IsSet("length") {
		snapshotLen = int32(ctx.Int("len"))
	}
	// Open device
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		misc.Log.Fatal(err)
	}
	defer handle.Close()

	// Set filter
	if ctx.IsSet("filter") {
		filter = ctx.String("filter")
	}
	err := handle.SetBPFFilter(filter)
	misc.Log.Infof("set SetBPFFilter: %v, err: %v", filter, err)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	ProcessPackets(packetSource.Packets())
}
