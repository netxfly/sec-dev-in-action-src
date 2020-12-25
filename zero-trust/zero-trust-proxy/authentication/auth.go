package authentication

import (
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/authentication/providers"
	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/logger"
	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/vars"
)

var ErrUnauthorized = errors.New("unauthorized request")

// JWTConfig JWT configuration
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type SecProxy struct {
	provider  providers.OAuth2Provider
	jwtConfig JWTConfig
}

// NewHeliosAuthentication creates a new authentication middleware instance
func NewSecProxyAuthentication(provider providers.OAuth2Provider, jwtSecret string, jwtExpiration time.Duration) SecProxy {
	return SecProxy{
		provider: provider,
		jwtConfig: JWTConfig{
			Secret:     jwtSecret,
			Expiration: jwtExpiration,
		},
	}
}

func (p SecProxy) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Debugf("Authenticating request %q", r.URL)
		if err := authenticate(p.jwtConfig.Secret, r); err != nil {
			logger.Log.Debugf("Authentication failed for %q", r.URL)
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			callback := scheme + "://" + r.Host + vars.CallbackPath

			url := p.provider.GetLoginURL(callback, r.RequestURI)

			logger.Log.Debugf("Redirecting to %s", url)

			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (p SecProxy) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("Handling callback request")
	encodedState := r.URL.Query().Get("state")
	state, err := base64.StdEncoding.DecodeString(encodedState)
	if err != nil {
		logger.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	profile, err := p.provider.FetchUser(r)
	if err != nil {
		logger.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Authorized. Redirecting to %s", string(state))

	expire := time.Now().Add(p.jwtConfig.Expiration)
	jwt, err := IssueJWTWithSecret(p.jwtConfig.Secret, profile.Email, expire)
	if err != nil {
		logger.Log.Error(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     vars.CookieName,
		Value:    jwt,
		Expires:  expire,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	http.Redirect(w, r, string(state), http.StatusFound)
}

func (p SecProxy) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     vars.CookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		Secure:   true,
		HttpOnly: true,
	})
}

func authenticate(jwtSecret string, r *http.Request) error {
	cookie, err := r.Cookie(vars.CookieName)
	token := r.Header.Get(vars.HeaderName)

	if err == http.ErrNoCookie && token == "" {
		return ErrUnauthorized
	}

	if token == "" && cookie != nil {
		token = cookie.Value
	}

	if !ValidateJWTWithSecret(jwtSecret, token) {
		return ErrUnauthorized
	}

	return nil
}
