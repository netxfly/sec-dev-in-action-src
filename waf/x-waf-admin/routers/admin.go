package routers

import (
	"fmt"
	"log"

	"github.com/go-macaron/captcha"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"sec-dev-in-action-src/waf/x-waf-admin/models"
	"sec-dev-in-action-src/waf/x-waf-admin/modules/util"
)

func Admin(ctx *macaron.Context, sess session.Store) {
	log.Println("session,", sess.Get("uid"))
	if sess.Get("uid") != nil {
		log.Println(sess.Get("uid"))
		sites, _ := models.ListSite()
		ctx.Data["sites"] = sites
		ctx.HTML(200, "AdminIndex")
	} else {
		ctx.Redirect("/login/")
	}
}

func SiteJSON(ctx *macaron.Context, sess session.Store) {
	if sess.Get("uid") != nil {
		log.Println(sess.Get("uid"))
		sites, _ := models.ListSite()
		ctx.JSON(200, sites)
	} else {
		ctx.Redirect("/login/")
	}
}

func Auth(ctx *macaron.Context, sess session.Store, cpt *captcha.Captcha) {
	if cpt.VerifyReq(ctx.Req) {
		username := ctx.Req.Form.Get("username")
		password := ctx.Req.Form.Get("password")
		ret, err := models.Auth(username, password)
		if err == nil && ret {
			sess.Set("uid", util.EncryptPass(fmt.Sprintf("%v+%v", username, util.EncryptPass(password))))
			ctx.Redirect("/admin/index/")
		} else {
			ctx.Redirect("/login/")
		}
	} else {
		message := "验证码输入错误"
		ctx.Data["message"] = message
		ctx.HTML(200, "error")
	}
}
