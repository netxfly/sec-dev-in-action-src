/*

Copyright (c) 2018 sec.lu

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

package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/elazarl/goproxy"

	"sec-dev-in-action-src/proxy-honeypot/agent/log"
	"sec-dev-in-action-src/proxy-honeypot/agent/util/api"
	"sec-dev-in-action-src/proxy-honeypot/agent/vars"
)

type (
	Meta struct {
		Req      *http.Request
		ReqBody  []byte
		ReqParam url.Values `json:"request_param"`
		Resp     *http.Response
		RespBody []byte
		Time     time.Time
		Session  int64
	}

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
)

func NewMeta(ctx *goproxy.ProxyCtx, reqParam url.Values, now time.Time) (meta *Meta) {
	meta = &Meta{Req: ctx.Req, Resp: ctx.Resp, Time: now, Session: ctx.Session, ReqParam: reqParam}
	return meta
}

func (m *Meta) readBody() {
	buf, err := ioutil.ReadAll(m.Req.Body)
	if err == nil {
		m.ReqBody = buf
	}
	// 用完恢复
	m.Req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	// 再用完再恢复
	// _ = m.Req.ParseForm()
	// m.ReqParam = m.Req.Form
	// m.Req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	respBuf, err := ioutil.ReadAll(m.Resp.Body)
	if err == nil {
		m.RespBody = respBuf
	}
	m.Resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBuf))
}

func (m *Meta) Parse() (record *HttpRecord) {
	record = &HttpRecord{}

	record.Session = m.Session
	record.Method = m.Req.Method
	record.RemoteAddr = m.Req.RemoteAddr
	record.StatusCode = m.Resp.StatusCode
	record.ContentLength = m.Resp.ContentLength
	record.Host = m.Resp.Request.Host
	// record.Port=m.Req
	record.Url = m.Resp.Request.URL.String()
	record.Scheme = m.Req.URL.Scheme
	record.Path = m.Req.URL.Path
	record.ReqHeader = m.Req.Header
	record.RespHeader = m.Resp.Header
	record.RequestParam = m.ReqParam
	record.RequestBody = m.ReqBody
	record.ResponseBody = m.RespBody
	record.VisitTime = m.Time

	return record
}

func (h *HttpRecord) print() {
	if vars.DebugMode {
		js, err := h.Json()
		log.Logger.Infof("data: %v, err: %v", string(js), err)
	}
}

func (h *HttpRecord) Json() (js []byte, err error) {
	js, err = json.Marshal(h)
	return js, err
}

func ReqHandlerFunc(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	vars.Cmap.Set(fmt.Sprintf("sess_%v", ctx.Session), req)
	if req != nil {
		buf, _ := ioutil.ReadAll(req.Body)
		reqTmp1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		req.Body = reqTmp1
		_ = req.ParseForm()
		params := req.Form
		reqTmp := ioutil.NopCloser(bytes.NewBuffer(buf))
		req.Body = reqTmp
		vars.Cmap.Set(fmt.Sprintf("sess_%v", ctx.Session), params)
	}
	return req, nil
}

func RespHandlerFunc(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if resp != nil {
		t, ok := vars.Cmap.Get(fmt.Sprintf("sess_%v", ctx.Session))
		defer vars.Cmap.Remove(fmt.Sprintf("sess_%v", ctx.Session))
		if ok {
			params, ok := t.(url.Values)
			if ok {
				meta := NewMeta(ctx, params, time.Now())
				meta.readBody()
				r := meta.Parse()
				r.print()
				data, err := r.Json()
				if err == nil {
					go func() {
						_ = api.Post(string(data))
						//log.HttpLogger.Info(data)
					}()
				}
			}
		}
	}

	return resp
}
