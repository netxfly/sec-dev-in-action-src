package main

import (
	"github.com/labstack/echo"

	"fmt"
	"html/template"
	"io"
	"net/http"

	"exchange_zero_trust_api/routers"
	"exchange_zero_trust_api/settings"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Static("/static", "public")
	e.Debug = true

	t := &Template{
		templates: template.Must(template.ParseGlob("./templates/*.html")),
	}
	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "working")
	})

	e.POST("/api/wbxml/", routers.DecodeWbxml) // 解析wbxml协议的接口
	e.POST("/api/send/", routers.SendNotice)   // 推送消息给用户的接口

	{
		// 移动端设备的授权路由
		e.GET("/active/:code", routers.Activation)
		e.POST("/a/activedevice", routers.ActiveDevice)
		e.POST("/a/ignoredevice", routers.IgnoreDevice)

		// PC端可信IP的授权路由
		e.GET("/pc/:code", routers.EwsActivation)
		e.POST("/a/active_ews", routers.ActiveEws)
		e.POST("/a/ignore_ews", routers.IgnoreEws)
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf("%v:%v", settings.Host, settings.Port)))
}
