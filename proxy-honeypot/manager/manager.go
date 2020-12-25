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

package main

import (
	"sec-dev-in-action-src/proxy-honeypot/manager/logger"
	"sec-dev-in-action-src/proxy-honeypot/manager/routers"
	"sec-dev-in-action-src/proxy-honeypot/manager/vars"

	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"fmt"
	"html/template"
	"net/http"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	m := macaron.Classic()
	m.Use(session.Sessioner())
	m.Use(csrf.Csrfer())
	m.Use(cache.Cacher())

	m.Use(macaron.Renderer(macaron.RenderOptions{
		Funcs: []template.FuncMap{map[string]interface{}{
			"Byte2Str": func(s []byte) string {
				return string(s)
			},
		}},
	}))

	m.Use(captcha.Captchaer(captcha.Options{
		URLPrefix:        "/captcha/",
		FieldIdName:      "captcha_id",
		FieldCaptchaName: "captcha",
		ChallengeNums:    4,
		Width:            150,
		Height:           40,
		Expiration:       600,
		CachePrefix:      "captcha_",
	}))

	m.Get("/", routers.Index)
	m.Get("/login/", routers.LoginIndex)
	m.Get("/logout/", routers.Logout)

	m.Group("/admin", func() {

		m.Get("/index/", routers.ListPassword)

		m.Get("/login/", routers.LoginIndex)
		m.Post("/login/", csrf.Validate, routers.Auth)

		m.Get("/dash/", routers.Dash)

		m.Group("/password/", func() {
			m.Get("", routers.ListPassword)
			m.Get("/list/", routers.ListPassword)
			m.Get("/list/:page", routers.ListPassword)
			m.Get("/list/:page", routers.ListPassword)

			m.Get("/site/:site/list/", routers.ListPasswordBySite)
			m.Get("/site/:site/list/:page", routers.ListPasswordBySite)

			m.Get("/detail/:id", routers.PasswordDetail)
		})

		m.Group("/record/", func() {
			m.Get("", routers.ListRecord)
			m.Get("/list/", routers.ListRecord)
			m.Get("/list/:page", routers.ListRecord)
			m.Get("/list/:page", routers.ListRecord)

			m.Get("/site/:site/list/", routers.ListRecordBySite)
			m.Get("/site/:site/list/:page", routers.ListRecordBySite)

			m.Get("/detail/:id", routers.RecordDetail)
		})

		m.Group("/tongji/", func() {
			m.Get("/site/", routers.TongjiPasswordBySite)
			m.Get("/site/:page", routers.TongjiPasswordBySite)

			m.Get("/urls/", routers.TongjiUrls)
		})

		m.Group("/user/", func() {
			m.Get("", routers.ListUser)
			m.Get("/list/", routers.ListUser)
			m.Get("/new/", routers.NewUser)
			m.Post("/new/", csrf.Validate, routers.DoNewUser)
			m.Get("/edit/:id", routers.EditUser)
			m.Post("edit/:id", csrf.Validate, routers.DoEditUser)
			m.Get("/del/:id", routers.DelUser)
		})

	})

	logger.Logger.Printf("Server is running on %s", fmt.Sprintf("%v:%v", vars.HttpHost, vars.HttpPort))
	logger.Logger.Println(http.ListenAndServe(fmt.Sprintf("%v:%v", vars.HttpHost, vars.HttpPort), m))
}
