package main

import (
	"fmt"
	"os"
	"runtime"

	"sec-dev-in-action-src/scanner/tcp-connect-scanner2/scanner"
	"sec-dev-in-action-src/scanner/tcp-connect-scanner2/util"
)

func main() {
	if len(os.Args) == 3 {
		ipList := os.Args[1]
		portList := os.Args[2]

		ips, err := util.GetIpList(ipList)
		ports, err := util.GetPorts(portList)
		_ = err

		tasks, _ := scanner.GenerateTask(ips, ports)
		scanner.RunTask(tasks)
		scanner.PrintResult()
	} else {
		fmt.Printf("%v iplist port\n", os.Args[0])
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
