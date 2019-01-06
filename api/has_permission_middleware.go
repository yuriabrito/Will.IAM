package api

import (
	"net/http"

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
			saID := mux.Vars(r)["id"]
			has, err := sasUC.HasPermission(saID, permission)
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
