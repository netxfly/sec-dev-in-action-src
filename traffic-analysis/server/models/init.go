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
	"fmt"

	"upper.io/db.v3"
	"upper.io/db.v3/mongo"

	"sec-dev-in-action-src/traffic-analysis/server/settings"
	"sec-dev-in-action-src/traffic-analysis/server/util"
)

var (
	DbConfig   DbCONF
	DbSettings db.ConnectionURL
	Session    db.Database
)

type DbCONF struct {
	DbType string
	DbHost string
	DbPort int64
	DbUser string
	DbPass string
	DbName string
}

func init() {
	cfg := settings.Cfg
	sec := cfg.Section("database")
	DbConfig.DbType = sec.Key("DB_TYPE").MustString("mongodb")
	DbConfig.DbHost = sec.Key("DB_HOST").MustString("127.0.0.1")
	DbConfig.DbPort = sec.Key("DB_PORT").MustInt64(27017)
	DbConfig.DbUser = sec.Key("DB_USER").MustString("user")
	DbConfig.DbPass = sec.Key("DB_PASS").MustString("password")
	DbConfig.DbName = sec.Key("DB_NAME").MustString("proxy_honeypot")

	_ = NewDbEngine()

}

func NewDbEngine() (err error) {
	switch DbConfig.DbType {
	case "mysql":
		util.Log.Info("will support mysql")
	case "mongodb":
		DbSettings = mongo.ConnectionURL{Host: fmt.Sprintf("%v:%v", DbConfig.DbHost, DbConfig.DbPort),
			User: DbConfig.DbUser, Password: DbConfig.DbPass, Database: DbConfig.DbName}
		Session, err = mongo.Open(DbSettings)
		util.Log.Warningf("settings: %v, session: %v, err: %v\n", DbSettings, Session, err)
		if err != nil {
			util.Log.Panicf("Connect Database failed, err: %v", err)
		}
		util.Log.Infof("DB Type: %v, Connect err status: %v", DbConfig.DbType, Session.Ping())
	}

	return err
}
