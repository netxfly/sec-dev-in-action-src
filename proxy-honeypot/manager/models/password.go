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
	"sec-dev-in-action-src/proxy-honeypot/manager/vars"

	"gopkg.in/mgo.v2/bson"

	"net/http"
	"net/url"
	"time"
)

type Password struct {
	Id                bson.ObjectId     `bson:"_id"`
	ResponseBody      string            `bson:"response_body"`
	RequestBody       string            `bson:"request_body"`
	DateStart         time.Time         `bson:"date_start"`
	URL               string            `bson:"url"`
	RequestParameters url.Values        `bson:"request_parameters"`
	FromIp            string            `bson:"from_ip"`
	Site              string            `bson:"site"`
	ResponseHeader    http.Header       `bson:"response_header"`
	RequestHeader     http.Header       `bson:"request_header"`
	Data              map[string]string `bson:"data"`
}

func ListPasswordByPage(page int) (passwords []Password, pages int, total int, err error) {

	coll := Session.DB(DataName).C("password")
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

	err = coll.Find(nil).Skip(i).Limit(vars.PageSize).All(&passwords)
	return passwords, pages, total, err
}

func ListPasswordBySite(site string, page int) (passwords []Password, pages int, total int, err error) {

	coll := Session.DB(DataName).C("password")
	total, _ = coll.Find(bson.M{"site": site}).Count()

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

	err = coll.Find(bson.M{"site": site}).Skip(i).Limit(vars.PageSize).All(&passwords)
	return passwords, pages, total, err
}

func PasswordDetail(id string) (Password, error) {
	var password Password
	coll := Session.DB(DataName).C("password")
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&password)
	return password, err
}
