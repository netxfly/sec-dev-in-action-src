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
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"sec-dev-in-action-src/proxy-honeypot/server/log"
	"sec-dev-in-action-src/proxy-honeypot/server/util"
)

type (
	HttpRecord struct {
		Id            int64       `json:"id"`
		Session       int64       `json:"session"`
		Method        string      `json:"method"`
		RemoteAddr    string      `json:"remote_addr"`
		StatusCode    int         `json:"status"`
		ContentLength int64       `json:"content_length"`
		Host          string      `json:"host"`
		Port          string      `json:"port"`
		Url           string      `json:"url"`
		Scheme        string      `json:"scheme"`
		Path          string      `json:"path"`
		ReqHeader     http.Header `json:"req_header"`
		RespHeader    http.Header `json:"resp_header"`
		RequestParam  url.Values  `json:"request_param"`
		RequestBody   []byte      `json:"request_body"`
		ResponseBody  []byte      `json:"response_body"`
		VisitTime     time.Time   `json:"visit_time"`
	}

	Record struct {
		Id                int64       `json:"id"`
		AgentIp           string      `json:"agent_ip"`
		AgentName         string      `json:"agent_name"`
		Remote            string      `json:"remote"`
		Method            string      `json:"method"`
		Status            int         `json:"status"`
		ContentLength     int64       `json:"content_length"`
		Host              string      `json:"host"`
		Port              string      `json:"port"`
		Url               string      `json:"url"`
		Scheme            string      `json:"scheme"`
		Path              string      `json:"path"`
		ReqHeader         http.Header `json:"req_header"`
		RespHeader        http.Header `json:"resp_header"`
		RequestBody       string      `json:"request_body"`
		ResponseBody      string      `json:"response_body" xorm:"LONGTEXT"`
		RequestParameters url.Values  `json:"request_parameters"`
		VisitTime         time.Time   `json:"visit_time"`
		Flag              int         `json:"flag"`
	}
)

func ParseHttpRecord(message string) (h HttpRecord, err error) {
	err = json.Unmarshal([]byte(message), &h)
	return h, err
}

func NewRecord(agentIp, agentName string, h HttpRecord) (record *Record) {

	return &Record{
		AgentIp:           agentIp,
		AgentName:         agentName,
		Remote:            util.Address2Ip(h.RemoteAddr),
		Method:            h.Method,
		Status:            h.StatusCode,
		ContentLength:     h.ContentLength,
		Host:              h.Host,
		Port:              h.Port,
		Url:               h.Url,
		Scheme:            h.Scheme,
		Path:              h.Path,
		ReqHeader:         h.ReqHeader,
		RespHeader:        h.RespHeader,
		RequestBody:       string(h.RequestBody),
		ResponseBody:      string(h.ResponseBody),
		RequestParameters: h.RequestParam,
		VisitTime:         h.VisitTime,
		Flag:              0,
	}
}

func (r *Record) Insert() (err error) {
	log.Logger.Warnf("remote: %v, url: %v", r.Remote, r.Url)
	if r.Remote != "" /*&& len(r.RequestParameters) > 0*/ {
		switch DbConfig.DbType {
		case "mysql":
			_, err = Engine.Table("record").Insert(r)
		case "mongodb":
			_, _ = GetSession()
			_, err = Session.Collection("record").Insert(r)
			log.Logger.Warnf("insert err: %v", err)
		}
	}
	return err
}
