package api

import (
	"net/http"

	"github.com/ghostec/Will.IAM/repositories"
)

func healthcheckHandler(
	repository repositories.Healthcheck,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := repository.Do(); err != nil {
			Write(w, http.StatusInternalServerError, `{"healthy": false}`)
		} else {
			Write(w, http.StatusOK, `{"healthy": true}`)
		}
	}
}
