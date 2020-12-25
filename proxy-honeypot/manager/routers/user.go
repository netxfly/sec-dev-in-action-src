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

package routers

import (
	"sec-dev-in-action-src/proxy-honeypot/manager/logger"
	"sec-dev-in-action-src/proxy-honeypot/manager/models"

	"github.com/go-macaron/captcha"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"

	"gopkg.in/macaron.v1"
)

func ListUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		users, _ := models.ListUser()
		ctx.Data["users"] = users
		ctx.Data["user"] = sess.Get("user")
		//log.Println(users)
		ctx.HTML(200, "user")
	} else {
		ctx.Redirect("/login/")
	}
}

func NewUser(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("admin") != nil {
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["user"] = sess.Get("user")
		ctx.HTML(200, "user_new")
	} else {
		ctx.Redirect("/login/")
	}
}

func DoNewUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		// ctx.Req.ParseForm()
		username := ctx.Req.Form.Get("username")
		password := ctx.Req.Form.Get("password")
		_ = models.NewUser(username, password)
		ctx.Redirect("/admin/user/")
	} else {
		ctx.Redirect("/login/")
	}
}

func EditUser(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("admin") != nil {
		Id := ctx.Params(":id")
		user, _ := models.GetUserById(Id)
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["user"] = user
		ctx.Data["username"] = user.UserName
		ctx.HTML(200, "user_edit")
	} else {
		ctx.Redirect("/login/")
	}
}

func DoEditUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		Id := ctx.Params(":id")
		username := ctx.Req.Form.Get("username")
		password := ctx.Req.Form.Get("password")
		_ = models.UpdateUser(Id, username, password)
		ctx.Redirect("/admin/user/")
	} else {
		ctx.Redirect("/login/")
	}
}

func DelUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		Id := ctx.Params(":id")
		_ = models.DelUser(Id)
		ctx.Redirect("/admin/user/")
	} else {
		ctx.Redirect("/login/")
	}
}

func Auth(ctx *macaron.Context, sess session.Store, cpt *captcha.Captcha) {
	if cpt.VerifyReq(ctx.Req) {
		username := ctx.Req.Form.Get("username")
		password := ctx.Req.Form.Get("password")
		ret, err := models.Auth(username, password)
		logger.Logger.Println(ret, err)
		if err == nil && ret {
			_ = sess.Set("admin", username)
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

func Logout(ctx *macaron.Context, sess session.Store) {
	_ = sess.Destory(ctx)
	ctx.Redirect("/login/")
}
