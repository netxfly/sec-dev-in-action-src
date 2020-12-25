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
	"sec-dev-in-action-src/proxy-honeypot/server/log"
	"sec-dev-in-action-src/proxy-honeypot/server/settings"

	"upper.io/db.v3"
	"upper.io/db.v3/mongo"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	"fmt"
)

var (
	DbConfig   DbCONF
	DbSettings db.ConnectionURL
	Session    db.Database

	Engine *xorm.Engine
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
	DbConfig.DbType = sec.Key("DB_TYPE").MustString("mysql")
	DbConfig.DbHost = sec.Key("DB_HOST").MustString("127.0.0.1")
	DbConfig.DbPort = sec.Key("DB_PORT").MustInt64(3306)
	DbConfig.DbUser = sec.Key("DB_USER").MustString("x-proxy")
	DbConfig.DbPass = sec.Key("DB_PASS").MustString("x@xsec.io")
	DbConfig.DbName = sec.Key("DB_NAME").MustString("x-proxy")

	_ = NewDbEngine()

}

func NewDbEngine() (err error) {
	switch DbConfig.DbType {
	case "mysql":
		dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8",
			DbConfig.DbUser, DbConfig.DbPass, DbConfig.DbHost, DbConfig.DbPort, DbConfig.DbName)
		Engine, err = xorm.NewEngine("mysql", dataSourceName)
		if err == nil {
			err = Engine.Ping()
			if err == nil {
				_ = Engine.Sync2(new(Record))
			}
		}

	case "mongodb":
		_, _ = GetSession()
	}

	return err
}

func GetSession() (db.Database, error) {
	var err error
	if Session == nil {
		DbSettings = mongo.ConnectionURL{Host: fmt.Sprintf("%v:%v", DbConfig.DbHost, DbConfig.DbPort), User: DbConfig.DbUser,
			Password: DbConfig.DbPass, Database: DbConfig.DbName}
		Session, err = mongo.Open(DbSettings)
		if err != nil {
			log.Logger.Panicf("Connect Database failed, err: %v", err)
		}
		Session.SetMaxOpenConns(100)
		log.Logger.Infof("DB Type: %v, DbSettings: %v, Connect err status: %v", DbConfig.DbType, DbSettings, Session.Ping())
	}

	return Session, err
}
