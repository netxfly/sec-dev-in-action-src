package main

import (
	"fmt"
	"net/url"

	"github.com/levigross/grequests"
)

func main() {
	proxyURL, err := url.Parse("http://sec.lu:8080") // Proxy URL
	if err != nil {
		panic(err)
	}

	resp, err := grequests.Get("http://mail.163.com/",
		&grequests.RequestOptions{Proxies: map[string]*url.URL{proxyURL.Scheme: proxyURL}})

	fmt.Printf("resp: %v, err: %v\n", resp, err)
}
