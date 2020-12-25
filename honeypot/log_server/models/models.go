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
	"log"
	"time"

	"gopkg.in/mgo.v2"

	"sec-dev-in-action-src/honeypot/log_server/settings"
	"sec-dev-in-action-src/honeypot/log_server/vars"
)

var (
	CollectionPacket  *mgo.Collection
	CollectionService *mgo.Collection
)

func ConnectMongoDb() (*mgo.Session, error) {
	Cfg := settings.Cfg
	sec := Cfg.Section("database")
	host := sec.Key("HOST").MustString("127.0.0.1")
	port := sec.Key("PORT").MustInt(27017)
	user := sec.Key("USER").MustString("honeypot")
	password := sec.Key("PASSWORD").MustString("xsec")
	database := sec.Key("DATABASE").MustString("honeypot")

	vars.MongodbName = database
	vars.CollPacket = sec.Key("COLL_PACKET").MustString("packet_info")
	vars.CollService = sec.Key("COLL_SERVICE").MustString("service")

	mongodbUrl := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v",
		user,
		password,
		host,
		port,
		database,
	)

	s, err := mgo.Dial(mongodbUrl)
	if err != nil {
		return nil, err

	}

	s.SetSafe(&mgo.Safe{})
	s.SetMode(mgo.Eventual, true)
	s.SetSocketTimeout(1 * time.Hour)

	return s, err
}

func init() {
	err := checkConnect()
	if err != nil {
		log.Panic(err)
	}
}

func checkConnect() error {
	var (
		err     error
		session *mgo.Session
	)

	if vars.Session == nil {
		session, err = ConnectMongoDb()
		if err == nil {
			vars.Session = session.Clone()
		}
	}

	CollectionService = vars.Session.DB(vars.MongodbName).C(vars.CollService)
	CollectionPacket = vars.Session.DB(vars.MongodbName).C(vars.CollPacket)

	return err
}
