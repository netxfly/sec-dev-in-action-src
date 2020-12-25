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
	"fmt"
	"time"

	"sec-dev-in-action-src/proxy-honeypot/manager/logger"
	"sec-dev-in-action-src/proxy-honeypot/manager/settings"

	"gopkg.in/mgo.v2"
)

var (
	Session  *mgo.Session
	Host     string
	Port     int
	USERNAME string
	PASSWORD string
	DataName string

	collAdmin *mgo.Collection
)

func init() {
	cfg := settings.Cfg
	sec := cfg.Section("MONGODB")
	Host = sec.Key("HOST").MustString("127.0.0.1")
	Port = sec.Key("PORT").MustInt(27017)
	USERNAME = sec.Key("USER").MustString("xproxy")
	PASSWORD = sec.Key("PASS").MustString("passw0rd")
	DataName = sec.Key("DATA").MustString("xproxy")
	err := NewMongodbClient()
	err = Session.Ping()
	logger.Logger.Infof("CONNECT MONGODB, err: %v", err)

	collAdmin = Session.DB(DataName).C("users")
	userCount, _ := collAdmin.Find(nil).Count()
	if userCount == 0 {
		_ = NewUser("xproxy", "x@xsec.io")
	}
}

// return a mongodb session
func NewMongodbClient() (err error) {
	url := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v", USERNAME, PASSWORD, Host, Port, DataName)
	Session, err = mgo.Dial(url)
	if err == nil {
		Session.SetSocketTimeout(1 * time.Hour)
	} else {
		logger.Logger.Panicf("connect mongodb failed, url: %v, err: %v", url, err)
	}
	return err
}
