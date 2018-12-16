package api

import (
	"fmt"
	"net/http"
	"net/url"

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
		code := qs["code"][0]
		authResult, err := provider.ExchangeCode(code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(qs["state"]) != 0 {
			v := url.Values{}
			v.Add("access_token", authResult.AccessToken)
			v.Add("email", authResult.Email)
			redirectTo := fmt.Sprintf("%s?%s", qs["state"][0], v.Encode())
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
