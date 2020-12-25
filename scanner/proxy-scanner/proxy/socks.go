package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"h12.io/socks"

	"sec-dev-in-action-src/scanner/proxy-scanner/models"
	"sec-dev-in-action-src/scanner/proxy-scanner/util"
)

var (
	SockProxyProtocol = map[string]int{"SOCKS4": socks.SOCKS4, "SOCKS4A": socks.SOCKS4A, "SOCKS5": socks.SOCKS5}
)

func CheckSockProxy(ip string, port int, protocol string) (isProxy bool, proxyInfo models.ProxyInfo, err error) {
	proxyInfo.Addr = ip
	proxyInfo.Port = port
	proxyInfo.Protocol = protocol

	proxy := fmt.Sprintf("%v:%v", ip, port)
	dialSocksProxy := socks.DialSocksProxy(SockProxyProtocol[protocol], proxy)
	tr := &http.Transport{Dial: dialSocksProxy}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(Timeout) * time.Second}

	resp, err := httpClient.Get(WebUrl)
	if err != nil {
		return false, proxyInfo, err
	}

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		// util.Log.Warningf("body: %v", string(body))
		if err != nil {
			return false, proxyInfo, err
		}
		if strings.Contains(string(body), "网易免费邮箱") {
			isProxy = true
		}
	}

	util.Log.Debugf("Checking proxy: %v, isProxy: %v", fmt.Sprintf("%v://%v:%v", protocol, ip, port), isProxy)

	return isProxy, proxyInfo, err
}
