package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device               = "en0"
	snapshotLength int32 = 1024
	promiscuous          = false
	timeout              = 30 * time.Second

	handle *pcap.Handle
	err    error
)

func main() {
	handle, err = pcap.OpenLive(device, snapshotLength, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		processPacket(packet)
		fmt.Println(strings.Repeat("-", 50))
	}
}

func processPacket(packet gopacket.Packet) {
	allLayer := packet.Layers()
	for _, layer := range allLayer {
		fmt.Printf("layer: %v\n", layer.LayerType())
	}

	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		fmt.Printf("Ethernet type: %v, source MAC: %v, destination MAC: %v\n",
			ethernetPacket.EthernetType, ethernetPacket.SrcMAC, ethernetPacket.DstMAC)
		fmt.Println(strings.Repeat("-", 50))
	}

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		fmt.Printf("proto: %v, from: %v, to: %v\n", ip.Protocol, ip.SrcIP, ip.DstIP)
	}

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		fmt.Printf("source port: %v, dest Port: %v\n", tcp.SrcPort, tcp.DstPort)
	}

	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		fmt.Printf("src port: %v, dst port: %v\n", udp.SrcPort, udp.DstPort)
	}

	appLayer := packet.ApplicationLayer()
	if appLayer != nil {
		fmt.Printf("application payload: %v\n", string(appLayer.Payload()))
	}

	err := packet.ErrorLayer()
	if err != nil {
		fmt.Printf("decode packet err: %v\n", err)
	}
}
