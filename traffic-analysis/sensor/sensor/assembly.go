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
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"

	"sec-dev-in-action-src/traffic-analysis/sensor/misc"
	"sec-dev-in-action-src/traffic-analysis/sensor/models"
)

type httpStreamFactory struct{}

type httpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

func (h *httpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hStream := &httpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	go hStream.run()
	return &hStream.r
}

func (h *httpStream) run() {
	buf := bufio.NewReader(&h.r)
	for {
		req, err := http.ReadRequest(buf)
		if err == io.EOF {
			return
		} else if err == nil {
			defer req.Body.Close()

			clientIp, dstIp := SplitNet2Ips(h.net)
			srcPort, dstPort := Transport2Ports(h.transport)

			httpReq, _ := models.NewHttpReq(req, clientIp, dstIp, dstPort)
			// send to sever
			go func(u string, req *models.HttpReq) {
				if !CheckSelfHtml(u, req) {
					misc.Log.Warnf("%v:%v -> %v(%v:%v), %v, %v, %v, req_param: %v, req_body: %v", httpReq.Client, srcPort, httpReq.Host, httpReq.Ip,
						httpReq.Port, httpReq.Method, httpReq.URL, httpReq.Header, httpReq.ReqParameters, httpReq.RequestBody)
					_ = SendHTML(req)
				}
			}(ApiUrl, httpReq)
		}
	}
}

func SplitNet2Ips(net gopacket.Flow) (client, host string) {
	ips := strings.Split(net.String(), "->")
	if len(ips) > 1 {
		client = ips[0]
		host = ips[1]
	}
	return client, host
}

func Transport2Ports(transport gopacket.Flow) (src, dst string) {
	ports := strings.Split(transport.String(), "->")
	if len(ports) > 1 {
		src = ports[0]
		dst = ports[1]
	}
	return src, dst
}

func CheckSelfHtml(ApiUrl string, req *models.HttpReq) (ret bool) {
	urlParsed, err := url.Parse(ApiUrl)
	if err == nil {
		apiIp := urlParsed.Host
		if apiIp == req.Host {
			ret = true
		}
		// misc.Log.Errorf("apiIp: %v, req.Host: %v, ret: %v", apiIp, req.Host, ret)
	}
	return ret
}

func ProcessPackets(packets chan gopacket.Packet) {
	streamFactory := &httpStreamFactory{}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	ticker := time.Tick(time.Minute)
	for {
		select {
		case packet := <-packets:
			if packet == nil {
				return
			}

			// 处理tcp/udp数据包
			processPacket(packet)
			// 处理DNS包
			parseDNS(packet)

			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				continue
			}

			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)

		case <-ticker:
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}
}

func SendHTML(req *models.HttpReq) error {
	reqJson, err := json.Marshal(req)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	urlApi := fmt.Sprintf("%v%v", ApiUrl, "/api/http/")
	secureKey := misc.MakeSign(timestamp, SecureKey)
	_, err = http.PostForm(urlApi, url.Values{"timestamp": {timestamp}, "secureKey": {secureKey}, "data": {string(reqJson)}})
	return err
}

func SendDns(dns *models.Dns) error {
	reqJson, err := json.Marshal(dns)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	urlApi := fmt.Sprintf("%v%v", ApiUrl, "/api/dns/")
	secureKey := misc.MakeSign(timestamp, SecureKey)
	_, err = http.PostForm(urlApi, url.Values{"timestamp": {timestamp}, "secureKey": {secureKey}, "data": {string(reqJson)}})
	return err
}
