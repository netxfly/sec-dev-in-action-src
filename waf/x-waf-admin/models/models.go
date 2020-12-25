package models

import (
	"fmt"
	"log"

	"sec-dev-in-action-src/waf/x-waf-admin/setting"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	Engine *xorm.Engine
	err    error
)

func init() {
	sec := setting.Cfg.Section("database")
	Engine, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		sec.Key("USER").String(),
		sec.Key("PASSWD").String(),
		sec.Key("HOST").String(),
		sec.Key("NAME").String()))
	if err != nil {
		log.Panicf("Faild to connect to database, err:%v", err)
	}

	// log.Println(Engine, err)
	Engine.Sync2(new(Site))
	Engine.Sync2(new(User))
	Engine.Sync2(new(Rules))

	ret, err := Engine.IsTableEmpty(new(User))
	if err == nil && ret {
		log.Printf("create new user:%v, password:%v\n", "admin", "x@xsec.io")
		NewUser("admin", "x@xsec.io")
	}

	ret, err = Engine.IsTableEmpty(new(Rules))
	if err == nil && ret {
		result, err := Engine.Exec(DefaultRules)
		log.Printf("Insert default waf rules, result: %v, err: %v\n", result, err)
	}
}
