package server

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// IsLogged looks for authentication information in the request in order to determine whether or not the user is logged
func (a *AuthServer) IsLogged(r *http.Request) (bool, jwt.Claims, map[string]interface{}) {
	var claims jwt.Claims
	var result map[string]interface{}
	log.Debug("isLogged")

	tokenString := a.Middleware.GetToken(r)
	if tokenString == "" {
		log.Debug("Did not find any token")
		return false, claims, result
	}

	parts := strings.Split(tokenString, " ")

	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		log.Info("failed-to-parse-cookie")
		return false, claims, result
	}

	parsed, err := jwt.ParseSigned(parts[1])
	if err != nil {
		log.Error("failed-to-parse-cookie-token", err)
		return false, claims, result
	}

	if err = parsed.Claims(&a.SigningKey.PublicKey, &claims, &result); err != nil {
		log.Error("failed-to-parse-claims ", err)
		return false, claims, result
	}

	if err = claims.Validate(jwt.Expected{Time: time.Now()}); err != nil {
		log.Error("failed-to-validate-claims ", err)
		return false, claims, result
	}

	return true, claims, result
}

// AuthenticationRequired is called for every admin request to determine if the request can be completed or shall be aborted
func (a *AuthServer) AuthenticationRequired(admin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		logged, _, user := a.IsLogged(c.Request)
		// TODO : verify claims against auths
		if !logged {
			redirectURL, err := url.Parse(c.Request.RequestURI)
			if err != nil {
				log.Error("Couldn't parse requestURI")
				c.JSON(http.StatusForbidden, gin.H{"error": "user needs to be admin to access this service"})
				c.Abort()
				return
			}
			redirectURL.Path = "/auth/login"
			params := redirectURL.Query()
			params.Set("redirect_uri", c.Request.Referer())
			redirectURL.RawQuery = params.Encode()
			c.Redirect(http.StatusSeeOther, redirectURL.String())
			c.Abort()
			return
		}
		if admin && user["is_admin"] != true {
			c.JSON(http.StatusForbidden, gin.H{"error": "user needs to be admin to access this service"})
			c.Abort()
			return
		}
		c.Next()
	}
}
