package google

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"

	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/authentication/providers"
)

type Provider struct {
	providers.OAuth2Provider
	oauth2 oauth2.Config
}

const (
	authEndpoint  = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenEndpoint = "https://oauth2.googleapis.com/token"
)

func NewGoogleProvider(config providers.OAuth2Config) providers.OAuth2Provider {
	return &Provider{
		oauth2: oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  authEndpoint,
				TokenURL: tokenEndpoint,
			},
			Scopes: []string{"email"},
		},
	}
}

func (provider Provider) FetchUser(r *http.Request) (providers.UserInfo, error) {
	var userInfo providers.UserInfo
	code := r.URL.Query().Get("code")

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	url := scheme + "://" + r.Host + r.URL.Path
	callback := oauth2.SetAuthURLParam("redirect_uri", url)

	// get access token
	token, err := provider.oauth2.Exchange(context.TODO(), code, callback)
	if err != nil {
		return userInfo, providers.ErrCodeExchange
	}

	// Parse the token
	jwtToken, _ := jwt.ParseWithClaims(token.Extra("id_token").(string), &providers.OIDClaims{}, nil)
	if jwtToken == nil {
		return userInfo, providers.ErrJWTParse
	}

	claims, ok := jwtToken.Claims.(*providers.OIDClaims)
	if !ok {
		return userInfo, providers.ErrJWTClaims
	}

	if claims.Email == "" {
		return userInfo, providers.ErrNoEmail
	}

	userInfo.Email = claims.Email

	return userInfo, nil
}

// GetLoginURL returns OAuth 2 login endpoint used to redirect users
func (provider Provider) GetLoginURL(callbackURL string, state string) string {
	s := base64.StdEncoding.EncodeToString([]byte(state))
	callback := oauth2.SetAuthURLParam("redirect_uri", callbackURL)
	return provider.oauth2.AuthCodeURL(s, callback)
}
