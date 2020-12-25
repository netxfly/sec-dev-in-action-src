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

package models

import (
	"gopkg.in/mgo.v2/bson"

	"net/http"
	"net/url"
	"time"

	"sec-dev-in-action-src/proxy-honeypot/manager/vars"
)

type (
	HttpRecord struct {
		Id            bson.ObjectId `bson:"_id"`
		Session       int64         `json:"session"`
		Method        string        `json:"method"`
		RemoteAddr    string        `json:"remote_addr" bson:"remote"`
		StatusCode    int           `json:"status"`
		ContentLength int64         `json:"content_length"`
		Host          string        `json:"host"`
		Port          string        `json:"port"`
		Url           string        `json:"url"`
		Scheme        string        `json:"scheme"`
		Path          string        `json:"path"`
		ReqHeader     http.Header   `json:"req_header"`
		RespHeader    http.Header   `json:"resp_header"`
		RequestParam  url.Values    `json:"request_param" bson:"requestparameters"`
		RequestBody   []byte        `json:"request_body"`
		ResponseBody  []byte        `json:"response_body"`
		VisitTime     time.Time     `json:"visit_time"`
	}
)

func ListRecordByPage(page int) (records []HttpRecord, pages int, total int, err error) {

	coll := Session.DB(DataName).C("record")
	total, _ = coll.Find(nil).Count()

	if int(total)%vars.PageSize == 0 {
		pages = int(total) / vars.PageSize
	} else {
		pages = int(total)/vars.PageSize + 1
	}

	if page >= pages {
		page = pages
	}

	if page < 1 {
		page = 1
	}

	i := (page - 1) * vars.PageSize
	if i < 0 {
		i = 0
	}

	err = coll.Find(nil).Skip(i).Limit(vars.PageSize).All(&records)
	return records, pages, total, err
}

func ListRecordBySite(site string, page int) (records []HttpRecord, pages int, total int, err error) {

	coll := Session.DB(DataName).C("record")
	total, _ = coll.Find(bson.M{"host": site}).Count()

	if int(total)%vars.PageSize == 0 {
		pages = int(total) / vars.PageSize
	} else {
		pages = int(total)/vars.PageSize + 1
	}

	if page >= pages {
		page = pages
	}

	if page < 1 {
		page = 1
	}

	i := (page - 1) * vars.PageSize
	if i < 0 {
		i = 0
	}

	err = coll.Find(bson.M{"host": site}).Skip(i).Limit(vars.PageSize).All(&records)
	return records, pages, total, err
}

func RecordDetail(id string) (HttpRecord, error) {
	var record HttpRecord
	coll := Session.DB(DataName).C("record")
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&record)
	return record, err
}
