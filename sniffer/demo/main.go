package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
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
		fmt.Println(packet.Dump())
	}
}
