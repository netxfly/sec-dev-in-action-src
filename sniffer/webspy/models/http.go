package models

import (
	"net/http"
	"net/url"
)

type HttpReq struct {
	Host          string
	Ip            string
	Client        string
	Port          string
	URL           *url.URL
	Header        http.Header
	RequestURI    string
	Method        string
	ReqParameters url.Values
}

func NewHttpReq(req *http.Request, client string, ip string, port string) (httpReq *HttpReq, err error) {
	err = req.ParseForm()
	return &HttpReq{Host: req.Host, Client: client, Ip: ip, Port: port, URL: req.URL, Header: req.Header,
		RequestURI: req.RequestURI, Method: req.Method, ReqParameters: req.Form}, err
}
