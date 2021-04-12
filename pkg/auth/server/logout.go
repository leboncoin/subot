package server

import (
	"net/http"
)

func (a *AuthServer) handleLogout(w http.ResponseWriter, _ *http.Request) {
	a.Middleware.UnsetToken(w)
	return
}
