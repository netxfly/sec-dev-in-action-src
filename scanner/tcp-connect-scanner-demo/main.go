package main

import (
	"fmt"
	"os"

	"sec-dev-in-action-src/scanner/tcp-connect-scanner-demo/scanner"
	"sec-dev-in-action-src/scanner/tcp-connect-scanner-demo/util"
)

func main() {
	if len(os.Args) == 3 {
		ipList := os.Args[1]
		portList := os.Args[2]
		ips, err := util.GetIpList(ipList)
		ports, err := util.GetPorts(portList)
		_ = err
		// fmt.Printf("iplist: %v, portList: %v, err: %v\n", ips, ports, err)
		for _, ip := range ips {
			for _, port := range ports {
				_, err := scanner.Connect(ip.String(), port)
				if err != nil {
					continue
				}
				fmt.Printf("ip: %v, port: %v is open\n", ip, port)
			}
		}

	} else {
		fmt.Printf("%v iplist port\n", os.Args[0])
	}
}
