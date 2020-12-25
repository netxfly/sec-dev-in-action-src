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
)

func DashPassword() (passwords []bson.M, err error) {
	coll := Session.DB(DataName).C("password")
	pipe := coll.Pipe([]bson.M{{"$group": bson.M{"_id": "$site", "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"count": -1}}, {"$limit": 20}})

	//resp := []bson.M{}

	err = pipe.All(&passwords)
	//if len(resp) > 0 {
	//	passwords = resp[0:20]
	//}
	return passwords, err
}

func DashUrls() (urls []bson.M, err error) {
	coll := Session.DB(DataName).C("urls")
	err = coll.Find(nil).Limit(20).All(&urls)

	return urls, err
}

func DashIps() (evilIps []bson.M, err error) {
	coll := Session.DB(DataName).C("evil_ips")
	err = coll.Find(nil).Limit(20).All(&evilIps)

	return evilIps, err
}

func DashTotal() (totalRecord int, totalPassword int, err error) {
	coll := Session.DB(DataName).C("proxy_honeypot")
	totalRecord, err = coll.Find(nil).Count()
	collPassword := Session.DB(DataName).C("password")
	totalPassword, err = collPassword.Find(nil).Count()
	return totalRecord, totalPassword, err
}
