/*

Copyright (c) 2017 xsec.io

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THEq
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/

package models

import (
	"net/http"
	"net/url"
	"time"
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

type EvilHttpReq struct {
	Id       int64
	Time     time.Time `bson:"time"`
	SensorIp string    `bson:"sensor_ip"`
	IsEvil   bool      `bson:"is_evil"`
	Data     HttpReq   `bson:"data"`
}

func NewEvilHttpReq(sensorIp string, isEvil bool, req HttpReq) (evilHttpReq *EvilHttpReq) {
	now := time.Now()
	return &EvilHttpReq{SensorIp: sensorIp, Time: now, IsEvil: isEvil, Data: req}
}

func (e *EvilHttpReq) Insert() {
	Session.Collection("http_req").Insert(e)
}

func ListEvilHttpReq() ([]EvilHttpReq, error) {
	evilHttpReqs := make([]EvilHttpReq, 0)
	res := Session.Collection("http_req").Find("-_id").OrderBy().Limit(500)
	err := res.All(&evilHttpReqs)
	return evilHttpReqs, err
}
