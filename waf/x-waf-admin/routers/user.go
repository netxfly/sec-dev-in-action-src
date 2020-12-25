package routers

import (
	"log"

	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"sec-dev-in-action-src/waf/x-waf-admin/models"
)

func ListUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("uid") != nil {
		users, _ := models.ListUser()
		ctx.Data["users"] = users
		ctx.HTML(200, "user")
	} else {
		ctx.Redirect("/login/")
	}
}

func UserJSON(ctx *macaron.Context, sess session.Store) {
	if sess.Get("uid") != nil {
		users, _ := models.ListUser()
		ctx.JSON(200, users)
	} else {
		ctx.Redirect("/login/")
	}
}

func NewUser(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("uid") != nil {
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.HTML(200, "newUser")
	} else {
		ctx.Redirect("/login/")
	}
}

func DoNewUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("uid") != nil {
		// ctx.Req.ParseForm()
		username := ctx.Req.Form.Get("username")
		password := ctx.Req.Form.Get("password")
		log.Println(username, password)
		models.NewUser(username, password)
		ctx.Redirect("/admin/user/")
	} else {
		ctx.Redirect("/login/")
	}
}

func EditUser(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("uid") != nil {
		Id := ctx.ParamsInt64(":id")
		user := models.User{Id: int64(Id)}
		models.Engine.Get(&user)
		log.Println(user)
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["user"] = user
		ctx.HTML(200, "EditUser")
	} else {
		ctx.Redirect("/login/")
	}
}

func DoEditUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("uid") != nil {
		// ctx.Req.ParseForm()
		Id := ctx.ParamsInt64(":id")
		username := ctx.Req.Form.Get("username")
		password := ctx.Req.Form.Get("password")
		models.UpdateUser(Id, username, password)
		ctx.Redirect("/admin/user/")
	} else {
		ctx.Redirect("/login/")
	}
}

func DelUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("uid") != nil {
		Id := ctx.ParamsInt64(":id")
		models.DelUser(Id)
		ctx.Redirect("/admin/user/")
	} else {
		ctx.Redirect("/login/")
	}
}
