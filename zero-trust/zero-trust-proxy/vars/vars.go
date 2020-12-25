package vars

import (
	"crypto/tls"

	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/config"
)

const (
	CookieName = "secProxy_Authorization"
	HeaderName = "secProxy-Jwt-Assertion"
)

var (
	ConfigPath   string
	ConfigFile   = "config.yaml"
	Conf         *config.Config
	TlsConfig    *tls.Config
	DebugMode    bool
	CurDir       string
	CaKey        string
	CaCert       string
	CallbackPath = "/.xsec/callback"
	LogoutPath   = "/.xsec/logout"
)
