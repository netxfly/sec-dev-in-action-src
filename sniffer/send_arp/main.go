package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	var ipAddr string
	var localAddr string
	var ip net.IP
	var localIp net.IP
	var eth string
	// Get a list of all interfaces.

	if len(os.Args) != 4 {
		fmt.Printf("usage: %v eth localip targetIp\n", os.Args[0])
		os.Exit(0)
	}

	eth = os.Args[1]
	localAddr = os.Args[2]
	ipAddr = os.Args[3]

	ip = net.ParseIP(ipAddr)
	localIp = net.ParseIP(localAddr)

	iFace, err := net.InterfaceByName(eth)
	if err != nil {
		panic(err)
	}

	err = send(iFace, localIp, ip)
	log.Printf("send arp to %v, iface: %v, err: %v\n", ip, iFace, err)

}

func send(iFace *net.Interface, localIp, ip net.IP) error {
	// Open up a pcap handle for packet reads/writes.
	handle, err := pcap.OpenLive(iFace.Name, 1024, true, pcap.BlockForever)
	if err != nil {
		return err
	}
	defer handle.Close()

	stop := make(chan struct{})
	go readARP(handle, iFace, stop)
	defer close(stop)
	for {
		if err := writeARP(handle, iFace, localIp, ip); err != nil {
			log.Printf("error writing packets on %v: %v", iFace.Name, err)
			return err
		}
		time.Sleep(1 * time.Second)
	}
}

func readARP(handle *pcap.Handle, iFace *net.Interface, stop chan struct{}) {
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()
	for {
		var packet gopacket.Packet
		select {
		case <-stop:
			return
		case packet = <-in:
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}
			arp := arpLayer.(*layers.ARP)
			log.Printf("arp.Operation: %v, srcIP %v, dstIp: %v, srcMac:%v, dstMac:%v\n",
				arp.Operation,
				net.IP(arp.SourceProtAddress),
				net.IP(arp.DstProtAddress),
				net.HardwareAddr(arp.SourceHwAddress),
				net.HardwareAddr(arp.DstHwAddress),
			)
		}
	}
}

func writeARP(handle *pcap.Handle, iFace *net.Interface, localIp, ip net.IP) error {
	eth := layers.Ethernet{
		SrcMAC:       iFace.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(iFace.HardwareAddr),
		SourceProtAddress: []byte(localIp),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	arp.DstProtAddress = ip
	_ = gopacket.SerializeLayers(buf, opts, &eth, &arp)
	if err := handle.WritePacketData(buf.Bytes()); err != nil {
		return err
	}
	return nil
}
