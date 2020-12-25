package proxy

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// ReverseProxyConfig configuration settings for a proxy instance
type ReverseProxyConfig struct {
	ConnectTimeout time.Duration
	Timeout        time.Duration
	IdleTimeout    time.Duration
}

// copy from: https://golang.org/src/net/http/httputil/reverseproxy.go
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// copy from: https://golang.org/src/net/http/httputil/reverseproxy.go
func NewSingleHostReverseProxy(target *url.URL, conf ReverseProxyConfig) http.Handler {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}

	return &httputil.ReverseProxy{
		FlushInterval: 200 * time.Millisecond,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				c, err := net.DialTimeout(network, addr, conf.ConnectTimeout)
				if err != nil {
					return c, err
				}

				if err := c.SetDeadline(time.Now().Add(conf.Timeout)); err != nil {
					return c, err
				}

				return c, err
			},
			TLSHandshakeTimeout:    10 * time.Second,
			IdleConnTimeout:        conf.IdleTimeout,
			MaxResponseHeaderBytes: 1 << 20,
			DisableCompression:     true,
		},
		Director: director,
	}
}
