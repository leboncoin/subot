package server

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"strings"
	"time"
)

func (a *AuthServer) handleLogin(w http.ResponseWriter, r *http.Request) {

	log.Debug("login")

	tokenString := a.Middleware.GetToken(r)
	if tokenString == "" {
		log.Debug("Did not find any token")
		a.NewLogin(w, r)
		return
	}

	redirectURI := r.FormValue("redirect_uri")
	if redirectURI == "" {
		redirectURI = "/"
	}

	parts := strings.Split(tokenString, " ")

	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		log.Info("failed-to-parse-cookie")
		a.NewLogin(w, r)
		return
	}

	parsed, err := jwt.ParseSigned(parts[1])
	if err != nil {
		log.Error("failed-to-parse-cookie-token", err)
		a.NewLogin(w, r)
		return
	}

	var claims jwt.Claims
	var result map[string]interface{}

	if err = parsed.Claims(&a.SigningKey.PublicKey, &claims, &result); err != nil {
		log.Error("failed-to-parse-claims ", err)
		a.NewLogin(w, r)
		return
	} else {
		log.Info("Successfully parsed claims....")
	}

	if err = claims.Validate(jwt.Expected{Time: time.Now()}); err != nil {
		log.Error("failed-to-validate-claims ", err)
		a.NewLogin(w, r)
		return
	}

	oauth2Token := &oauth2.Token{
		TokenType:   parts[0],
		AccessToken: parts[1],
		Expiry:      claims.Expiry.Time(),
	}

	token := oauth2Token.WithExtra(map[string]interface{}{
		"csrf": result["csrf"],
	})

	a.Redirect(w, r, token, redirectURI)
}

func (a *AuthServer) NewLogin(w http.ResponseWriter, r *http.Request) {
	redirectURI := r.FormValue("redirect_uri")

	if redirectURI == "" {
		redirectURI = "/"
	}

	scopes := []string{"openid", "profile", "email", "federated:id", "groups"}

	stateToken := encode(&StateToken{
		RedirectURI: redirectURI,
		Entropy:     RandomString(),
	})

	http.SetCookie(w, &http.Cookie{
		Name:    stateCookieName,
		Value:   stateToken,
		Path:    "/",
		Expires: time.Now().Add(time.Hour),
	})

	authCodeURL := a.Oauth2Config(scopes).AuthCodeURL(stateToken)
	log.Debug("redirect to ", authCodeURL)

	http.Redirect(w, r, authCodeURL, http.StatusTemporaryRedirect)
}