package server

import (
	"fmt"
	"github.com/coreos/go-oidc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

func (a *AuthServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	log.Debug("entered auth callback")
	var (
		err   error
		token *oauth2.Token
		stateToken string
	)

	ctx := oidc.ClientContext(r.Context(), a.Client)
	oauth2Config := a.Oauth2Config(nil)

	log.Debug("List cookies : ", r.Cookies())
	cookieState, err := r.Cookie(stateCookieName)
	if err != nil {
		log.Error("failed-to-fetch-cookie-state", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Authorization redirect callback from OAuth2 auth flow.
	if errMsg := r.FormValue("error"); errMsg != "" {
		http.Error(w, errMsg+": "+r.FormValue("error_description"), http.StatusBadRequest)
		return
	}

	if stateToken = cookieState.Value; stateToken != r.FormValue("state") {
		log.Error("failed-with-unexpected-state-token", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		http.Error(w, fmt.Sprintf("no code in request: %q", r.Form), http.StatusBadRequest)
		return
	}

	token, err = oauth2Config.Exchange(ctx, code)
	if err != nil {
		log.Error("failed-to-fetch-dex-token", err)
		switch e := err.(type) {
		case *oauth2.RetrieveError:
			http.Error(w, string(e.Body), e.Response.StatusCode)
			return
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	log.Debug("Reset state cookie")
	http.SetCookie(w, &http.Cookie{
		Name:   stateCookieName,
		Path:   "/",
		MaxAge: -1,
	})

	claims, err := a.Verifier.Verify(r.Context(), token)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to verify ID token: %v", err), http.StatusInternalServerError)
		return
	}

	authToken, err := a.Issuer.Issue(claims)
	if err != nil {
		log.Error("failed-to-issue-token", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	a.Redirect(w, r, authToken, decode(stateToken).RedirectURI)
}
