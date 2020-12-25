package routers

import (
	"github.com/go-macaron/csrf"

	"gopkg.in/macaron.v1"
)

func Index(ctx *macaron.Context) {
	ctx.Redirect("/login/")
	//return 403, "forbidden by nginx"
	//ctx.HTML(200, "index")
}

func LoginIndex(ctx *macaron.Context, x csrf.CSRF) {
	ctx.Data["csrf_token"] = x.GetToken()
	ctx.HTML(200, "login")
}
