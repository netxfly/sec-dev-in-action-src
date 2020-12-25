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
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/toolkits/slice"
	"strings"
	"sync"
	"time"

	"sec-dev-in-action-src/honeypot/agent/logger"
	"sec-dev-in-action-src/honeypot/agent/models"
)

var (
	mux sync.Mutex
)

func ProcessPacket(packets chan gopacket.Packet) {
	for packet := range packets {
		processPacket(packet)
	}
}

func processPacket(packet gopacket.Packet) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, ok := ipLayer.(*layers.IPv4)
		if ok {
			switch ip.Protocol {
			case layers.IPProtocolTCP:
				tcpLayer := packet.Layer(layers.LayerTypeTCP)
				if tcpLayer != nil {
					tcp, _ := tcpLayer.(*layers.TCP)

					srcPort := SplitPortService(tcp.SrcPort.String())
					dstPort := SplitPortService(tcp.DstPort.String())
					isHttp := false

					applicationLayer := packet.ApplicationLayer()
					if applicationLayer != nil {
						// Search for a string inside the payload
						if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
							isHttp = true
						}
					}

					connInfo := models.NewConnectionInfo("tcp", ip.SrcIP.String(), srcPort, ip.DstIP.String(), dstPort, isHttp)

					go func(info *models.ConnectionInfo, ) {
						if !IsInWhite(info) &&
							!CheckSelfPacker(info) &&
							(tcp.SYN && !tcp.ACK) {
							err := SendPacker(info)
							logger.Log.Debugf("[TCP] %v:%v -> %v:%v, err: %v", ip.SrcIP, tcp.SrcPort.String(),
								ip.DstIP, tcp.DstPort.String(), err)
						}
					}(connInfo)
				}
			}
		}
	}
}

func SendPacker(connInfo *models.ConnectionInfo) (err error) {
	packetInfo := models.NewPacketInfo(connInfo, time.Now())
	jsonPacket, err := packetInfo.String()
	if err != nil {
		return err
	}

	go logger.LogReport.WithField("api", "/api/packet/").Info(jsonPacket)

	return err
}

func CheckSelfPacker(p *models.ConnectionInfo) (ret bool) {
	if slice.ContainsString(SensorIps, p.SrcIp) || p.DstIp == ApiIp || p.SrcIp == ApiIp {
		ret = true
	}
	return ret
}

func SplitPortService(portService string) (port string) {
	t := strings.Split(portService, "(")
	if len(t) > 0 {
		port = t[0]
	}
	return port
}
