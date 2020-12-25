package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sec-dev-in-action-src/scanner/proxy-scanner/models"
	"sec-dev-in-action-src/scanner/proxy-scanner/util"
)

var (
	HttpProxyProtocol = []string{"http", "https"}
	WebUrl            = "http://email.163.com/"
)

func CheckHttpProxy(ip string, port int, protocol string) (isProxy bool, proxyInfo models.ProxyInfo, err error) {
	proxyInfo.Addr = ip
	proxyInfo.Port = port
	proxyInfo.Protocol = protocol

	rawProxyUrl := fmt.Sprintf("%v://%v:%v", protocol, ip, port)
	proxyUrl, err := url.Parse(rawProxyUrl)
	if err != nil {
		return false, proxyInfo, err
	}

	Transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	client := &http.Client{Transport: Transport, Timeout: time.Duration(Timeout) * time.Second}

	resp, err := client.Get(WebUrl)
	if err != nil {
		return false, proxyInfo, err
	}

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		// util.Log.Warningf("body: %v", string(body))
		if err != nil {
			return false, proxyInfo, err
		}

		if strings.Contains(string(body), "<title>网易免费邮箱") {
			isProxy = true
		}
	}

	util.Log.Debugf("Checking proxy: %v, isProxy: %v", rawProxyUrl, isProxy)

	return isProxy, proxyInfo, err
}
