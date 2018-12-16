package api

import (
	"net/http"

	"github.com/ghostec/Will.IAM/oauth2"
)

func authenticationBuildURLHandler(
	provider oauth2.Provider,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Write(w, http.StatusOK, provider.BuildAuthURL("http://localhost:3001/authentication/sso_test"))
	}
}

func authenticationHandler(
	provider oauth2.Provider,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		qs := r.URL.Query()
		if len(qs["code"]) == 0 {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		// TODO: ExchangeCode
		if len(qs["state"]) != 0 {
			redirectTo := qs["state"][0]
			http.Redirect(w, r, redirectTo, http.StatusSeeOther)
			return
		}
		Write(w, http.StatusOK, "")
	}
}

func authenticationSSOTestHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Write(w, http.StatusOK, "")
	}
}
