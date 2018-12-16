package api

import (
	"net/http"

	"github.com/ghostec/Will.IAM/oauth2"
)

func authenticationBuildURLHandler(
	provider oauth2.Provider,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Write(w, http.StatusOK, provider.BuildAuthURL(""))
	}
}
