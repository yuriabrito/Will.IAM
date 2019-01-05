package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ghostec/Will.IAM/oauth2"
	"github.com/ghostec/Will.IAM/usecases"
)

func authenticationBuildURLHandler(
	provider oauth2.Provider,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		qs := r.URL.Query()
		if len(qs["origin"]) == 0 {
			Write(
				w, http.StatusUnprocessableEntity,
				`{ "error": "querystrings.origin is required" }`,
			)
			return
		}
		authURL := provider.BuildAuthURL(qs["origin"][0])
		http.Redirect(w, r, authURL, http.StatusSeeOther)
	}
}

func authenticationExchangeCodeHandler(
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
			v.Add("accessToken", authResult.AccessToken)
			v.Add("email", authResult.Email)
			v.Add("origin", qs["state"][0])
			redirectTo := fmt.Sprintf("/sso?%s", v.Encode())
			http.Redirect(w, r, redirectTo, http.StatusSeeOther)
			return
		}
		Write(w, http.StatusOK, "")
	}
}

func authenticationValidHandler(
	provider oauth2.Provider, sasUC usecases.ServiceAccounts,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		qs := r.URL.Query()
		if len(qs["origin"]) == 0 {
			Write(
				w, http.StatusUnprocessableEntity,
				`{ "error": "querystrings.origin is required" }`,
			)
			return
		}
		if len(qs["accessToken"]) == 0 {
			Write(
				w, http.StatusUnprocessableEntity,
				`{ "error": "querystrings.accessToken is required" }`,
			)
			return
		}
		authResult, err := sasUC.AuthenticateAccessToken(qs["accessToken"][0])
		if err != nil {
			// TODO: check if err is non-authorized
			v := url.Values{}
			v.Add("origin", qs["origin"][0])
			http.Redirect(
				w, r, fmt.Sprintf("/sso/auth/do?%s", v.Encode()), http.StatusSeeOther,
			)
			return
		}
		v := url.Values{}
		v.Add("origin", qs["origin"][0])
		v.Add("accessToken", authResult.AccessToken)
		http.Redirect(
			w, r, fmt.Sprintf("%s?%s", qs["origin"][0], v.Encode()),
			http.StatusSeeOther,
		)
	}
}
