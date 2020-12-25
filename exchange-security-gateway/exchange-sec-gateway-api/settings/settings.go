package settings

import (
	"gopkg.in/ini.v1"

	"exchange_zero_trust_api/logger"
)

var (
	Cfg  *ini.File
	Host string
	Port int
)

func init() {
	var err error
	source := "conf/app.ini"
	Cfg, err = ini.Load(source)

	if err != nil {
		logger.Log.Panicln(err)
	}

	Host = Cfg.Section("").Key("HOST").MustString("127.0.0.1")
	Port = Cfg.Section("").Key("PORT").MustInt(5000)
}
