package util

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type ProxyAddr struct {
	IP   string
	Port int
}

func ReadProxyAddr(fileName string) (sliceProxyAddr []ProxyAddr) {
	proxyFile, err := os.Open(fileName)
	if err != nil {
		Log.Fatalf("Open proxy file err, %v", err)
	}

	defer proxyFile.Close()

	scanner := bufio.NewScanner(proxyFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ipPort := strings.TrimSpace(scanner.Text())

		if ipPort == "" {
			continue
		}

		t := strings.Split(ipPort, ":")
		ip := t[0]
		port, err := strconv.Atoi(t[1])
		if err == nil {
			proxyAddr := ProxyAddr{IP: ip, Port: port}
			sliceProxyAddr = append(sliceProxyAddr, proxyAddr)
		}
	}

	return sliceProxyAddr
}

