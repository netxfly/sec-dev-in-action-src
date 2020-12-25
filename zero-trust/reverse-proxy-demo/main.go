package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
	UpstreamUrl string
}

func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remote, err := url.Parse(p.UpstreamUrl)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

func main() {
	addr := ":8888"
	proxyHandle := &ReverseProxy{UpstreamUrl: "http://127.0.0.1:8000"}
	log.Printf("proxy addr: %v, Upstream: %v\n", addr, proxyHandle)
	err := http.ListenAndServe(addr, proxyHandle)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
