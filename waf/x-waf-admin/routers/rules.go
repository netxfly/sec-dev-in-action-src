package routers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"sec-dev-in-action-src/waf/x-waf-admin/models"
	"sec-dev-in-action-src/waf/x-waf-admin/modules/util"
	"sec-dev-in-action-src/waf/x-waf-admin/setting"

	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

func ListRules(ctx *macaron.Context, sess session.Store) {
	log.Println("session,", sess.Get("uid"))
	if sess.Get("uid") != nil {
		log.Println(sess.Get("uid"))
		rules, _ := models.ListAllRules()
		ctx.Data["RuleInfo"] = models.RuleInfo
		ctx.Data["Rules"] = rules
		ctx.HTML(200, "rules")
	} else {
		ctx.Redirect("/login/")
	}
}

func NewRule(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("uid") != "" {
		ctx.Data["csrf_token"] = x.GetToken()
		ruleType := ctx.Params("type")
		ctx.Data["ruleType"] = ruleType
		ctx.HTML(200, "newRule")
	} else {
		ctx.Redirect("/login/")
	}
}

func DoNewRule(ctx *macaron.Context, sess session.Store) {
	if sess.Get("uid") != nil {
		log.Println(sess.Get("uid"))
		ruleType := ctx.Params("type")
		rule := ctx.Req.Form.Get("rule")
		log.Println(ruleType, rule)
		models.NewRule(ruleType, rule)
		ctx.Redirect("/admin/rule/")
	} else {
		ctx.Redirect("/login/")
	}
}

func EditRule(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("uid") != nil {
		Id := ctx.ParamsInt64(":id")
		rule := models.Rules{Id: int64(Id)}
		models.Engine.Get(&rule)
		log.Println(rule)
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["rule"] = rule
		ctx.HTML(200, "EditRule")
	} else {
		ctx.Redirect("/login/")
	}
}

func DoEditRule(ctx *macaron.Context, sess session.Store) {
	if sess.Get("uid") != nil {
		// ctx.Req.ParseForm()
		Id := ctx.ParamsInt64(":id")
		rule := ctx.Req.Form.Get("rule")
		models.EditRule(Id, rule)
		ctx.Redirect("/admin/rule/")
	} else {
		ctx.Redirect("/login/")
	}
}

func DelRule(ctx *macaron.Context, sess session.Store) {
	if sess.Get("uid") != nil {
		Id := ctx.ParamsInt64(":id")
		models.DelRule(Id)
		ctx.Redirect("/admin/rule/")
	} else {
		ctx.Redirect("/login/")
	}
}

func SyncRule(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") != nil {
		timestamp := time.Now().Unix()
		hash := util.MakeMd5(setting.AppKey + util.MakeMd5(fmt.Sprintf("%v%v", timestamp, setting.AppKey)))
		for _, server := range setting.APIServers {
			server = strings.TrimSpace(server)
			url := fmt.Sprintf("http://%s:%v/api/rule/sync/?hash=%v&timestamp=%v", server, setting.HTTPPort, hash, timestamp)
			log.Println(url)
			resp, err := http.Get(url)
			if err == nil {
				body, err := ioutil.ReadAll(resp.Body)
				log.Println(string(body), err)
				flash.Success(string(body))
			} else {
				flash.Success(err.Error())
			}
		}
		ctx.Redirect("/admin/rule/")
	} else {
		ctx.Redirect("/login/")
	}
}

func SyncRuleApi(ctx *macaron.Context) {
	timestamp := ctx.Query("timestamp")
	hash := ctx.Query("hash")
	if util.MakeMd5(setting.AppKey+util.MakeMd5(fmt.Sprintf("%v%v", timestamp, setting.AppKey))) == hash {
		rules, _ := models.ListAllRules()
		for k, item := range rules {
			ruleFile := fmt.Sprintf("%v/%v.rule", setting.RulePath, k)
			log.Println(ruleFile)
			file, err := os.Create(ruleFile)
			if err == nil {
				ruleJson, err := json.Marshal(item)
				log.Println(string(ruleJson), err)
				file.WriteString(string(ruleJson))
			}
			file.Close()
		}

		if util.ReloadNginx() == nil {
			ret := util.Message{Status: 0, Message: "successful"}
			ctx.JSON(200, &ret)
		} else {
			ret := util.Message{Status: 1, Message: "reload nginx configure faild"}
			ctx.JSON(200, &ret)
		}
	} else {
		ret := util.Message{Status: 2, Message: "invalid hash parameter"}
		ctx.JSON(200, &ret)
	}
}
