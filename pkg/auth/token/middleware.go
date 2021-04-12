package token

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

// Middleware interface
type Middleware interface {
	SetToken(http.ResponseWriter, string, time.Time) error
	UnsetToken(http.ResponseWriter)
	GetToken(*http.Request) string
}

type middleware struct {
	secureCookies bool
}

// NewMiddleware returns a new instance of the Middleware interface
func NewMiddleware(secureCookies bool) Middleware {
	return &middleware{secureCookies: secureCookies}
}

const authCookieName = "auth_token"
const stateCookieName = "state_token"
// NumCookies is the number of cookie to split the token into
const NumCookies = 15
const maxCookieSize = 4000

// UnsetToken invalidates all the tokens related to our application (used for logouts)
func (m *middleware) UnsetToken(w http.ResponseWriter) {
	for i := 0; i < NumCookies; i++ {
		http.SetCookie(w, &http.Cookie{
			Name:     authCookieName + strconv.Itoa(i),
			Path:     "/",
			MaxAge:   -1,
			//Secure:   m.secureCookies,
			//HttpOnly: true,
			Secure:  false,
		})
	}
}

// SetToken writes the token string to a token that is added to the response
func (m *middleware) SetToken(w http.ResponseWriter, tokenStr string, expiry time.Time) error {
	tokenLength := len(tokenStr)
	if tokenLength > maxCookieSize*NumCookies {
		return errors.New("token is too long to fit in cookies")
	}

	for i := 0; i < NumCookies; i++ {
		if len(tokenStr) > maxCookieSize {
			http.SetCookie(w, &http.Cookie{
				Name:     authCookieName + strconv.Itoa(i),
				Value:    tokenStr[:maxCookieSize],
				Path:     "/",
				Expires:  expiry,
				//HttpOnly: true,
				//Secure:   m.secureCookies,
				Secure:  false,
			})
			tokenStr = tokenStr[maxCookieSize:]
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:     authCookieName + strconv.Itoa(i),
				Value:    tokenStr,
				Path:     "/",
				Expires:  expiry,
				//HttpOnly: true,
				//Secure:   m.secureCookies,
				Secure:   false,
			})
			break
		}
	}
	return nil
}

// GetToken reads and parses the token from the request cookies
func (m *middleware) GetToken(r *http.Request) string {
	log.Debug("Get token")
	authCookie := ""
	for i := 0; i < NumCookies; i++ {
		cookie, err := r.Cookie(authCookieName + strconv.Itoa(i))
		if err == nil {
			authCookie += cookie.Value
		}
	}
	return authCookie
}
