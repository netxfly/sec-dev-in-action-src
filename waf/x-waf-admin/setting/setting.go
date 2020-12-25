package setting

import (
	"log"
	"strings"

	"gopkg.in/ini.v1"
	"gopkg.in/macaron.v1"
)

var (
	// Application settings
	AppVer   string
	ProdMode bool

	// server settings
	HTTPPort    int
	AppKey      string
	NginxBin    string
	NginxVhosts string
	APIServers  []string
	RulePath    string

	Cfg *ini.File
)

func init() {
	log.SetPrefix("[xsec-waf]")
	var err error
	source := "conf/app.ini"
	Cfg, err = ini.Load(source)
	// log.Println(Cfg, err)
	if err != nil {
		log.Panicln(err)
	}

	if Cfg.Section("").Key("RUN_MODE").String() == "prod" {
		ProdMode = true
		macaron.Env = macaron.PROD
		macaron.ColorLog = false
	}

	sec := Cfg.Section("server")
	HTTPPort = sec.Key("HTTP_PORT").MustInt(5000)
	AppKey = sec.Key("API_KEY").MustString("www.xsec.io")
	NginxBin = sec.Key("NGINX_BIN").MustString("/usr/local/openresty/nginx/sbin/nginx")
	NginxVhosts = sec.Key("NGINX_VHOSTS").MustString("/usr/local/openresty/nginx/conf/vhosts/")
	APIServers = strings.Split(sec.Key("API_SERVERS").MustString("127.0.0.1"), ",")

	RulePath = Cfg.Section("waf").Key("RULE_PATH").MustString("/usr/local/openresty/nginx/conf/waf/rules/")
}
