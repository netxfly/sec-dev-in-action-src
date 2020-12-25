package util

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/authentication"
	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/authentication/providers"
	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/authentication/providers/google"
	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/authorization"

	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/logger"
	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/proxy"
	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/vars"
)

func init() {
	vars.CurDir, _ = GetCurDir()

	vars.CaKey = filepath.Join(vars.CurDir, "./certs/sever.key")
	vars.CaCert = filepath.Join(vars.CurDir, "./certs/server.cert")
}

func GetCurDir() (string, error) {
	dir, err := filepath.Abs(filepath.Dir("./"))
	if err != nil {
		return "", err
	}
	return dir, err
}

func SetupRouter() *mux.Router {
	// mux router对象
	router := mux.NewRouter()
	// 上游服务器的map，key为upstream名称，value为proxy.NewSingleHostReverseProxy的返回值，类型为http.Handler，
	upstreams := make(map[string]http.Handler)
	// oauth2的配置
	oauth2conf := providers.OAuth2Config{
		ClientID:     vars.Conf.Identity.ClientID,
		ClientSecret: vars.Conf.Identity.ClientSecret,
		AuthURL:      vars.Conf.Identity.OAuth2.AuthURL,
		TokenURL:     vars.Conf.Identity.OAuth2.TokenURL,
		ProfileURL:   vars.Conf.Identity.OAuth2.ProfileURL,
	}
	// oauth2提供者的配置信息
	var provider providers.OAuth2Provider
	switch vars.Conf.Identity.Provider {
	case "google":
		provider = google.NewGoogleProvider(oauth2conf)
	default:
		logger.Log.Fatalf("%q provider is not supported", vars.Conf.Identity.Provider)
	}
	// 创建一个认证对象，传入的参数是oauth2的provider，jwt的密钥与超时时间。
	authN := authentication.NewSecProxyAuthentication(provider, vars.Conf.JWT.Secret, vars.Conf.JWT.Expires)
	// 给路由：/.xsec/callback设置处理器：authN.CallbackHandler
	router.PathPrefix(vars.CallbackPath).HandlerFunc(authN.CallbackHandler)
	// 给路由：/.xsec/logout设置处理器：authN.LogoutHandler
	router.PathPrefix(vars.LogoutPath).HandlerFunc(authN.LogoutHandler)

	// 枚举upstream信息，生成一个map[string]http.Handler, key为upstream名称，value为proxy.NewSingleHostReverseProxy的返回值
	for _, upstream := range vars.Conf.Upstreams {
		upstreamURL, err := url.Parse(upstream.URL)
		logger.Log.Debugf("upstreamUrl: %v, err: %v", upstreamURL, err)
		if err != nil {
			logger.Log.Fatalf("Cannot parse upstream %q URL: %v", upstream.Name, err)
		}

		proxyConf := proxy.ReverseProxyConfig{
			ConnectTimeout: upstream.ConnectTimeout,
			IdleTimeout:    vars.Conf.Server.IdleTimeout,
			Timeout:        vars.Conf.Server.Timeout,
		}
		reverseProxy := proxy.NewSingleHostReverseProxy(upstreamURL, proxyConf)
		upstreams[upstream.Name] = reverseProxy
	}

	// 枚举全部的路由信息
	for _, route := range vars.Conf.Routes {
		// 只匹配host为route.Host的子路由，如router.Host("www.example.com")，router.Host("{subdomain}.domain.com")等
		h := router.Host(route.Host).Subrouter()

		for _, path := range route.HTTP.Paths {
			// 通过path.Upstream中指定的名字，获取到相应的upstream
			upstream := upstreams[path.Upstream]
			// 如果upstream不存在，说明配置不正确，直接退出程序
			if upstream == nil {
				logger.Log.Fatalf("Upstream %q for route %q not found", path.Upstream, route.Host)
				break
			}
			// 为Host为route.Host的请求指定后端stream
			// h.PathPrefix(path.Path).Handler(upstream)

			authZ := authorization.NewAuthorization(route.Rules)
			logger.Log.Debugf("path.Authentication: %v", path.Authentication)

			if path.Authentication {
				h.PathPrefix(path.Path).Handler(authN.Middleware(authZ.Middleware(upstream)))
			} else {
				h.PathPrefix(path.Path).Handler(authZ.Middleware(upstream))
			}

		}
	}

	return router
}

func Start(c *cli.Context) error {
	if c.IsSet("debug") {
		vars.DebugMode = c.Bool("debug")
	}
	if c.IsSet("config") {
		vars.ConfigFile = c.String("config") + ".yaml"
	}

	vars.ConfigPath = filepath.Join(vars.CurDir, "conf", vars.ConfigFile)

	if vars.DebugMode {
		logger.Log.Logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.Log.Logger.SetLevel(logrus.InfoLevel)
	}

	cb, err := ioutil.ReadFile(vars.ConfigPath)
	if err != nil {
		logger.Log.Fatalf("Error loading configuration: %v", err)
	}

	if err = yaml.Unmarshal(cb, &vars.Conf); err != nil {
		logger.Log.Fatalf("Error parsing configuration: %v", err)
	}

	vars.TlsConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
	}

	address := fmt.Sprintf("%v:%v", vars.Conf.Server.ListenIP, vars.Conf.Server.ListenPort)
	server := &http.Server{
		Addr:           address,
		WriteTimeout:   vars.Conf.Server.Timeout,
		ReadTimeout:    vars.Conf.Server.Timeout,
		IdleTimeout:    vars.Conf.Server.IdleTimeout,
		TLSConfig:      vars.TlsConfig,
		MaxHeaderBytes: 1 << 20, // 1mb
		Handler:        SetupRouter(),
	}

	logger.Log.Infof("debug_mode: %v, config_path: %v, addr: %v", vars.DebugMode, vars.ConfigFile, address)

	err = server.ListenAndServeTLS(vars.Conf.Server.TLSContext.CertificatePath, vars.Conf.Server.TLSContext.PrivateKeyPath)
	if err != nil {
		logger.Log.Errorf("start server failed, err: %v", err)
	}

	return err
}
