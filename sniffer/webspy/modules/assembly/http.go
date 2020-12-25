package assembly

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"

	"sec-dev-in-action-src/sniffer/webspy/logger"
	"sec-dev-in-action-src/sniffer/webspy/models"
	"sec-dev-in-action-src/sniffer/webspy/vars"
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

			defer func() {
				_ = req.Body.Close()
			}()

			clientIp, dstIp := SplitNet2Ips(h.net)
			srcPort, dstPort := Transport2Ports(h.transport)

			httpReq, _ := models.NewHttpReq(req, clientIp, dstIp, dstPort)
			logger.Log.Infof("httpReq: %v", httpReq)

			// send to sever
			go func(addr string, req *models.HttpReq, ) {
				reqInfo := fmt.Sprintf("%v:%v -> %v(%v:%v), %v, %v, %v, %v", httpReq.Client, srcPort, httpReq.Host, httpReq.Ip,
					httpReq.Port, httpReq.Method, httpReq.URL, httpReq.Header, httpReq.ReqParameters)
				logger.Log.Infof("reqInfo: %v", reqInfo)

				SendHTML(reqInfo)
				//if !CheckSelfHtml(addr, req) {
				//	SendHTML(req)
				//}
			}(vars.HttpHost, httpReq)
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

func CheckSelfHtml(host string, req *models.HttpReq) (ret bool) {
	if host == strings.Split(req.Host, ":")[0] {
		ret = true
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

			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)

		case <-ticker:
			assembler.FlushOlderThan(time.Now().Add(time.Second * -20))
		}
	}
}

func SendHTML(reqInfo string) {
	vars.Data.Put(reqInfo)
}
