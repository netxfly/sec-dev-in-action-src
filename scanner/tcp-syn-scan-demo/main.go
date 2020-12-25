// Copyright 2012 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

// synscan implements a TCP syn scanner on top of pcap.
// It's more complicated than arpscan, since it has to handle sending packets
// outside the local network, requiring some routing and ARP work.
//
// Since this is just an example program, it aims for simplicity over
// performance.  It doesn't handle sending packets very quickly, it scans IPs
// serially instead of in parallel, and uses gopacket.Packet instead of
// gopacket.DecodingLayerParser for packet processing.  We also make use of very
// simple timeout logic with time.Since.
//
// Making it blazingly fast is left as an exercise to the reader.
package main

import (
	"fmt"
	"os"

	"sec-dev-in-action-src/scanner/tcp-syn-scan-demo/scan"
	"sec-dev-in-action-src/scanner/tcp-syn-scan-demo/util"
)

func main() {
	if len(os.Args) == 3 {
		util.CheckRoot()
		
		ipList := os.Args[1]
		portList := os.Args[2]
		ips, err := util.GetIpList(ipList)
		ports, err := util.GetPorts(portList)
		_ = err

		for _, ip := range ips {
			for _, port := range ports {
				ip1, port1, err1 := scan.SynScan(ip.String(), port)
				if err1 == nil && port1 > 0 {
					fmt.Printf("%v:%v is open\n", ip1, port1)
				}
			}
		}
	} else {
		fmt.Printf("%v iplist port\n", os.Args[0])
	}
}
