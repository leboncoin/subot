package server

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"strings"
	"time"
)

func (a *AuthServer) userInfo(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"header": r.Header.Get("Authorization")}).Debug("User info")

	tokenString := a.Middleware.GetToken(r)
	if tokenString == "" {
		log.Error("Did not find any token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	parts := strings.Split(tokenString, " ")

	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		log.Error("Not enouth parts, or part is not bearer ")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	parsed, err := jwt.ParseSigned(parts[1])
	if err != nil {
		log.Error("failed-to-parse-authorization-token ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var claims jwt.Claims
	var userInfo map[string]interface{}

	if err = parsed.Claims(&a.SigningKey.PublicKey, &claims, &userInfo); err != nil {
		log.Error("failed-to-parse-claims ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = claims.Validate(jwt.Expected{Time: time.Now()}); err != nil {
		log.Error("failed-to-validate-claims", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(userInfo)
	if err != nil {
		log.Error("error while encoding userinfo", err)
	}
}
