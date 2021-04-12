package auth

import (
	"crypto/x509"
	"encoding/pem"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"github.com/leboncoin/subot/pkg/auth/dex"
	"github.com/leboncoin/subot/pkg/auth/server"
	"github.com/leboncoin/subot/pkg/auth/token"
	"time"
)

// NewServer returns a http handler for dex and authent servers as well as the auth server
func NewServer() (http.Handler, server.AuthServer) {
	baseURL := viper.GetString("analytics_url")
	issuerURL := baseURL + "/dex/issuer"
	redirectURL := baseURL + "/auth/callback"

	privPem, _ := pem.Decode([]byte(viper.GetString("dex_private_key")))
	if privPem == nil {
		log.Error("Failed to parse pem")
	}
	parsedKey, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if err != nil {
		log.Error("Could not parse key", err)
	}

	tokenVerifier := token.NewVerifier(viper.GetString("dex_client_id"), issuerURL)
	tokenIssuer := token.NewIssuer(token.NewGenerator(parsedKey), 12*time.Hour)
	tokenMiddleware := token.NewMiddleware(true)

	authServer := server.NewAuthServer(
		baseURL,
		viper.GetString("dex_secret"),
		*tokenVerifier,
		tokenMiddleware,
		tokenIssuer,
		parsedKey,
	)
	authHandler := authServer.NewAuthHandler()

	dexServer, err := dex.NewDexServer(issuerURL, redirectURL)
	if err != nil {
		return nil, authServer
	}

	handler := http.NewServeMux()
	handler.Handle("/dex/issuer/", dexServer)
	handler.Handle("/auth/", authHandler)

	return handler, authServer
}
