package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"sec-dev-in-action-src/waf/x-waf-admin/routers"
	"sec-dev-in-action-src/waf/x-waf-admin/setting"

	"fmt"
	"strings"
)

const APP_VER = "0.1"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	setting.AppVer = APP_VER
}

func main() {
	m := macaron.Classic()

	m.Use(macaron.Renderer())
	m.Use(session.Sessioner())
	m.Use(csrf.Csrfer())
	m.Use(cache.Cacher())

	m.Use(captcha.Captchaer(captcha.Options{
		URLPrefix:        "/captcha/",
		FieldIdName:      "captcha_id",
		FieldCaptchaName: "captcha",
		ChallengeNums:    6,
		Width:            220,
		Height:           50,
		Expiration:       600,
		CachePrefix:      "captcha_",
	}))

	m.Get("/", routers.Index)
	m.Get("/login/", routers.LoginIndex)
	m.Post("/login/", csrf.Validate, routers.Auth)

	m.Group("/admin", func() {
		m.Get("/index/", routers.Admin)
		m.Group("/site/", func() {
			m.Get("", routers.Admin)
			m.Get("/list/", routers.Admin)
			m.Get("/new/", routers.NewSite)
			m.Post("/new/", csrf.Validate, routers.DoNewSite)
			m.Get("/edit/:id", routers.EditSite)
			m.Post("/edit/:id", csrf.Validate, routers.DoEditSite)
			m.Get("/del/:id", routers.DelSite)
			m.Get("/sync/", routers.SyncSite)
			m.Get("/sync/:id", routers.SyncSiteById)
			m.Get("/json/", routers.SiteJSON)
		})

		m.Group("/rule/", func() {
			m.Get("", routers.ListRules)
			m.Get("/list/", routers.ListRules)
			m.Get("/new/:type", routers.NewRule)
			m.Post("/new/:type", csrf.Validate, routers.DoNewRule)
			m.Get("/edit/:id", routers.EditRule)
			m.Post("/edit/:id", csrf.Validate, routers.DoEditRule)
			m.Get("/del/:id", routers.DelRule)
			m.Get("/sync/", routers.SyncRule)
		})

		m.Group("/user/", func() {
			m.Get("", routers.ListUser)
			m.Get("/list/", routers.ListUser)
			m.Get("/new/", routers.NewUser)
			m.Post("/new/", csrf.Validate, routers.DoNewUser)
			m.Get("/edit/:id", routers.EditUser)
			m.Post("edit/:id", csrf.Validate, routers.DoEditUser)
			m.Get("/del/:id", routers.DelUser)
			m.Get("/json/", routers.UserJSON)
		})
	})

	m.Group("/api", func() {
		m.Get("/site/sync/", routers.SyncSiteApi)
		m.Get("/rule/sync/", routers.SyncRuleApi)
	})

	log.Printf("xsec waf admin %s", setting.AppVer)
	log.Printf("Run mode %s", strings.Title(macaron.Env))
	log.Printf("Server is running on %s", fmt.Sprintf("0.0.0.0:%v", setting.HTTPPort))
	log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", setting.HTTPPort), m))
}
