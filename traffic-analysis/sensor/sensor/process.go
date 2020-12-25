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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"sec-dev-in-action-src/traffic-analysis/sensor/misc"
	"sec-dev-in-action-src/traffic-analysis/sensor/models"
	"sec-dev-in-action-src/traffic-analysis/sensor/settings"
)

func processPacket(packet gopacket.Packet) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		if ip != nil {
			switch ip.Protocol {
			case layers.IPProtocolTCP:
				tcpLayer := packet.Layer(layers.LayerTypeTCP)
				if tcpLayer != nil {
					tcp, _ := tcpLayer.(*layers.TCP)

					srcPort := tcp.SrcPort.String()
					dstPort := tcp.DstPort.String()
					connInfo := models.NewConnectionInfo("tcp", ip.SrcIP.String(), srcPort, ip.DstIP.String(), dstPort)

					go func(u string, info *models.ConnectionInfo) {
						if tcp.SYN && !tcp.ACK && !CheckSelfPacker(ApiUrl, connInfo) {
							misc.Log.Debugf("[TCP] %v:%v -> %v:%v, syn: %v, ack: %v, seq: %v, ack: %v, psh: %v", ip.SrcIP, tcp.SrcPort.String(),
								ip.DstIP, tcp.DstPort.String(), tcp.SYN, tcp.ACK, tcp.Seq, tcp.Ack, tcp.PSH)
							_ = SendPacker(info)
						}
					}(ApiUrl, connInfo)
				}
				//
				//case layers.IPProtocolUDP:
				//	udpLayer := packet.Layer(layers.LayerTypeUDP)
				//	if udpLayer != nil {
				//		udp, _ := udpLayer.(*layers.UDP)
				//
				//		srcPort := udp.SrcPort.String()
				//		dstPort := udp.DstPort.String()
				//		connInfo := models.NewConnectionInfo("tcp", ip.SrcIP.String(), srcPort, ip.DstIP.String(), dstPort)
				//
				//		go func(u string, info *models.ConnectionInfo) {
				//			if !CheckSelfPacker(u, info) {
				//				misc.Log.Debugf("[UDP] %v:%v -> %v:%v", ip.SrcIP, udp.SrcPort.String(), ip.DstIP, udp.DstPort.String())
				//				SendPacker(info)
				//			}
				//		}(ApiUrl, connInfo)
				//	}
			}
		}
	}
}

func parseDNS(packet gopacket.Packet) {
	var eth layers.Ethernet
	var ip4 layers.IPv4
	var udp layers.UDP
	var dns layers.DNS
	parser := gopacket.NewDecodingLayerParser(
		layers.LayerTypeEthernet, &eth, &ip4, &udp, &dns)
	decodedLayers := make([]gopacket.LayerType, 0)
	err := parser.DecodeLayers(packet.Data(), &decodedLayers)
	if err != nil {
		return
	}
	srcIp := ip4.SrcIP
	dstIp := ip4.DstIP
	for _, v := range dns.Questions {
		dns := models.NewDns(srcIp.String(), dstIp.String(), v.Type.String(), string(v.Name))
		go func(dns *models.Dns) {
			misc.Log.Debugf("%v -> %v, dns_type: %v, dns_name: %v", srcIp, dstIp, v.Type, string(v.Name))
			_ = SendDns(dns)
		}(dns)
	}
}

func SendPacker(connInfo *models.ConnectionInfo) (err error) {
	infoJson, err := json.Marshal(connInfo)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	urlApi := fmt.Sprintf("%v%v", ApiUrl, "/api/packet/")
	secureKey := misc.MakeSign(timestamp, SecureKey)

	resp, err := http.PostForm(urlApi, url.Values{"timestamp": {timestamp}, "secureKey": {secureKey}, "data": {string(infoJson)}})
	_ = resp
	return err
}

func CheckSelfPacker(ApiUrl string, p *models.ConnectionInfo) (ret bool) {
	urlParsed, err := url.Parse(ApiUrl)
	if err == nil {
		apiHost := urlParsed.Host
		apiIp := strings.Split(apiHost, ":")[0]
		sensorIp := settings.Ips[0]

		if p.SrcIp == sensorIp && p.DstIp == apiIp || p.SrcIp == apiIp && p.DstIp == sensorIp {
			ret = true
		}
		// misc.Log.Errorf("srcIp:%v, sensorIp: %v, DstIp: %v, ApiSeverIp: %v, ret: %v", p.SrcIp, sensorIp, p.DstIp, apiIp, ret)
	}
	return ret
}
