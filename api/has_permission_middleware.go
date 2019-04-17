package api

import (
	"net/http"
	"strings"

	"github.com/ghostec/Will.IAM/usecases"
	"github.com/gorilla/mux"
	"github.com/topfreegames/extensions/middleware"
)

func hasPermissionMiddlewareBuilder(
	sasUC usecases.ServiceAccounts,
) func(string, http.Handler) http.Handler {
	return func(permission string, next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := middleware.GetLogger(r.Context())
			saID, ok := getServiceAccountID(r.Context())
			if !ok {
				l.Error("No ServiceAccountID in r.Context()")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if id, ok := mux.Vars(r)["id"]; ok {
				permission = strings.Replace(permission, "{id}", id, -1)
			}
			has, err := sasUC.WithContext(r.Context()).
				HasPermissionString(saID, permission)
			if err != nil {
				l.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !has {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
