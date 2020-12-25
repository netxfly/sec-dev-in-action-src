package routers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"sec-dev-in-action-src/waf/x-waf-admin/models"
	"sec-dev-in-action-src/waf/x-waf-admin/modules/util"
	"sec-dev-in-action-src/waf/x-waf-admin/setting"
)

func NewSite(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("uid") != "" {
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.HTML(200, "newSite")
	} else {
		ctx.Redirect("/login/")
	}
}

func DoNewSite(ctx *macaron.Context, sess session.Store) {
	if sess.Get("uid") != nil {
		log.Println(sess.Get("uid"))
		siteName := ctx.Req.Form.Get("sitename")
		port := ctx.Req.Form.Get("port")
		Port, _ := strconv.Atoi(port)

		// get nginx upstream backend address
		backaddr := ctx.Req.Form.Get("backendaddr")
		backendaddr := strings.Split(backaddr, "\r\n")
		BackendAddr := make([]string, 0)
		for _, v := range backendaddr {
			if v == "" {
				continue
			}
			v = strings.TrimSpace(v)
			BackendAddr = append(BackendAddr, v)
		}

		// get nginx upstream unreal address
		unrealAddr := ctx.Req.Form.Get("unreal_addr")
		unrealAddrs := strings.Split(unrealAddr, "\r\n")
		UnrealAddr := make([]string, 0)
		for _, v := range unrealAddrs {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			UnrealAddr = append(UnrealAddr, v)
		}

		ssl := ctx.Req.Form.Get("ssl")
		debugLevel := ctx.Req.Form.Get("debuglevel")

		log.Println(siteName, BackendAddr, ssl, debugLevel)
		models.NewSite(siteName, Port, BackendAddr, UnrealAddr, ssl, debugLevel)
		ctx.Redirect("/admin/site/list/")
	} else {
		ctx.Redirect("/login/")
	}
}

func EditSite(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("uid") != nil {
		Id := ctx.ParamsInt64(":id")
		site := models.Site{Id: int64(Id)}
		models.Engine.Get(&site)
		log.Println(site)
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["site"] = site
		ctx.HTML(200, "EditSite")
	} else {
		ctx.Redirect("/login/")
	}
}

func DoEditSite(ctx *macaron.Context, sess session.Store) {
	if sess.Get("uid") != nil {
		// ctx.Req.ParseForm()
		Id := ctx.ParamsInt64(":id")
		siteName := ctx.Req.Form.Get("sitename")
		port := ctx.Req.Form.Get("port")
		Port, _ := strconv.Atoi(port)

		// get nginx upstream backend address
		backaddr := ctx.Req.Form.Get("backendaddr")
		backendaddr := strings.Split(backaddr, "\r\n")
		BackendAddr := make([]string, 0)
		for _, v := range backendaddr {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			BackendAddr = append(BackendAddr, v)
		}

		// get nginx upstream unreal address
		unrealAddr := ctx.Req.Form.Get("unreal_addr")
		unrealAddrs := strings.Split(unrealAddr, "\r\n")
		UnrealAddr := make([]string, 0)
		for _, v := range unrealAddrs {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			UnrealAddr = append(UnrealAddr, v)
		}

		ssl := ctx.Req.Form.Get("ssl")
		debugLevel := ctx.Req.Form.Get("debuglevel")
		log.Println(Id, siteName, BackendAddr, ssl, debugLevel)
		log.Println(models.UpdateSite(Id, siteName, Port, BackendAddr, UnrealAddr, ssl, debugLevel))
		ctx.Redirect("/admin/site/list/")
	} else {
		ctx.Redirect("/login/")
	}
}

func DelSite(ctx *macaron.Context, sess session.Store) {
	if sess.Get("uid") != nil {
		Id := ctx.ParamsInt64(":id")
		models.DelSite(Id)
		ctx.Redirect("/admin/index/")
	} else {
		ctx.Redirect("/login/")
	}
}

func SyncSite(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") != nil {
		timestamp := time.Now().Unix()
		hash := util.MakeMd5(setting.AppKey + util.MakeMd5(fmt.Sprintf("%v%v", timestamp, setting.AppKey)))
		for _, server := range setting.APIServers {
			server = strings.TrimSpace(server)
			url := fmt.Sprintf("http://%s:%v/api/site/sync/?hash=%v&timestamp=%v", server, setting.HTTPPort, hash, timestamp)
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
		ctx.Redirect("/admin/site/list/")
	} else {
		ctx.Redirect("/login/")
	}
}

func SyncSiteById(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") != nil {
		Id := ctx.ParamsInt64(":id")
		timestamp := time.Now().Unix()
		hash := util.MakeMd5(setting.AppKey + util.MakeMd5(fmt.Sprintf("%v%v", timestamp, setting.AppKey)))
		for _, server := range setting.APIServers {
			server = strings.TrimSpace(server)
			url := fmt.Sprintf("http://%s:%v/api/site/sync/?id=%v&hash=%v&timestamp=%v", server, setting.HTTPPort, Id, hash, timestamp)
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
		ctx.Redirect("/admin/site/list/")
	} else {
		ctx.Redirect("/login/")
	}
}

func SyncSiteApi(ctx *macaron.Context) {
	timestamp := ctx.Query("timestamp")
	hash := ctx.Query("hash")
	id := ctx.Query("id")
	if util.MakeMd5(setting.AppKey+util.MakeMd5(fmt.Sprintf("%v%v", timestamp, setting.AppKey))) == hash {
		Id, err := strconv.Atoi(id)
		log.Println(Id, err)
		var sites []models.Site
		if err == nil {
			sites, err = models.ListSiteById(int64(Id))
		} else {
			sites, err = models.ListSite()
		}
		log.Println(sites, err)
		for _, site := range sites {
			ctx.Data["site"] = site
			proxyConfig, err := ctx.HTMLString("proxy", ctx.Data)
			log.Println(proxyConfig, err)
			util.WriteNginxConf(proxyConfig, site.SiteName, setting.NginxVhosts)

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
