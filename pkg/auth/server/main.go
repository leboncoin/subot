package server

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"github.com/leboncoin/subot/pkg/auth/token"
	"strings"
)

const stateCookieName = "state_token"

// AuthServer structure
type AuthServer struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	SigningKey   *rsa.PrivateKey

	Verifier   token.Verifier
	Issuer     token.Issuer
	Middleware token.Middleware

	// Does the provider use "offline_access" scope to request a refresh token
	// or does it use "access_type=offline" (e.g. Google)?
	OfflineAsScope bool

	Client *http.Client
}

// Oauth2Config returns and oauth2.Config object from the given scopes
func (a *AuthServer) Oauth2Config(scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     a.ClientID,
		ClientSecret: a.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   strings.TrimRight(viper.GetString("analytics_url") + "/dex/issuer", "/") + "/auth",
			TokenURL:  strings.TrimRight(viper.GetString("analytics_url") + "/dex/issuer", "/") + "/token",
			AuthStyle: oauth2.AuthStyleInHeader,
		},
		Scopes:      scopes,
		RedirectURL: a.RedirectURI,
	}
}

// StateToken represents a token transiting through the login process
type StateToken struct {
	RedirectURI string `json:"redirect_uri"`
	Entropy     string `json:"entropy"`
}

// RandomString returns a random 32 characters string
func RandomString() string {
	s := make([]byte, 32)
	rand.Read(s)
	return hex.EncodeToString(s)
}

func decode(raw string) *StateToken {
	data, _ := base64.StdEncoding.DecodeString(raw)

	var token *StateToken
	err := json.Unmarshal(data, &token)
	if err != nil {
		return token
	}
	return token
}

func encode(token *StateToken) string {
	j, _ := json.Marshal(token)

	return base64.StdEncoding.EncodeToString(j)
}

// Redirect to the redirectURI when specified
func (a *AuthServer) Redirect(w http.ResponseWriter, r *http.Request, token *oauth2.Token, redirectURI string) {
	log.Debug("redirect")

	redirectURL, err := url.Parse(redirectURI)
	if err != nil {
		log.Error("failed-to-parse-redirect-url", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	csrfToken, ok := token.Extra("csrf").(string)
	if !ok {
		log.Error("failed-to-extract-csrf-token", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.Middleware.SetToken(w, token.TokenType+" "+token.AccessToken, token.Expiry)
	if err != nil {
		log.Error("invalid-token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//if redirectURL.Host != os.Getenv("API_URL") && redirectURL.Host != os.Getenv("FRONT_URL") && redirectURL.Host != "" {
	//	log.Error("invalid-redirect", fmt.Errorf("Unsupported redirect uri: %s", redirectURI))
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	params := redirectURL.Query()
	params.Set("csrf_token", csrfToken)
	redirectURL.RawQuery = params.Encode()

	http.Redirect(w, r, redirectURL.String(), http.StatusTemporaryRedirect)
}

// NewAuthServer returns a new instance of the authServer from the given parameters
func NewAuthServer(apiURL string, secret string, verifier token.Verifier, middleware token.Middleware, issuer token.Issuer, signingKey *rsa.PrivateKey) AuthServer {
	var a AuthServer

	a.SigningKey = signingKey
	a.Client = http.DefaultClient
	a.ClientID = "support-analytics"
	a.ClientSecret = secret
	a.RedirectURI = fmt.Sprintf("%s/auth/callback", apiURL)
	a.Verifier = verifier
	a.Middleware = middleware
	a.Issuer = issuer

	return a
}

// NewAuthHandler returns a custom http handler dedicated to authentication paths
func (a *AuthServer) NewAuthHandler() http.Handler {
	handler := http.NewServeMux()
	handler.HandleFunc("/auth/login", a.handleLogin)
	handler.HandleFunc("/auth/logout", a.handleLogout)
	handler.HandleFunc("/auth/callback", a.handleCallback)
	handler.HandleFunc("/auth/userinfo", a.userInfo)

	log.Debug("Returning handler : ", handler)

	return handler
}
