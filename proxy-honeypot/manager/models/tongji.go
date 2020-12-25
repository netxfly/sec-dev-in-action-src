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

	"sec-dev-in-action-src/proxy-honeypot/manager/vars"
)

func TongjiPasswordBySite(page int) (passwords []bson.M, pages int, total int, err error) {
	coll := Session.DB(DataName).C("password")
	pipe := coll.Pipe([]bson.M{{"$group": bson.M{"_id": "$site", "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"count": -1}}})

	resp := []bson.M{}
	err = pipe.All(&resp)
	total = len(resp)

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
	// logger.Logger.Warnf("i: %v, page*vars.PageSize: %v, resp: %v", i, page*vars.PageSize, resp)
	if page*vars.PageSize-i > len(resp) {
		passwords = resp
	} else {
		passwords = resp[i : page*vars.PageSize]
	}
	return passwords, pages, total, err
}

func TongjiUrls(page int) (urls []bson.M, pages int, total int, err error) {
	coll := Session.DB(DataName).C("proxy_honeypot")
	pipe := coll.Pipe([]bson.M{{"$group": bson.M{"_id": "$host", "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"count": -1}}, {"allowDiskUse": true}})

	resp := []bson.M{}
	err = pipe.All(&resp)
	total = len(resp)

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

	urls = resp[i : page*vars.PageSize]
	return urls, pages, total, err
}
