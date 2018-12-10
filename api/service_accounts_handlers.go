package api

import (
	"net/http"

	"github.com/ghostec/Will.IAM/usecases"
	"github.com/gorilla/mux"
)

func serviceAccountsHasPermissionHandler(
	serviceAccountsUseCase usecases.ServiceAccounts,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		qs := r.URL.Query()
		permissionSl := qs["permission"]
		if len(permissionSl) == 0 {
			Write(w, http.StatusBadRequest, `{"error": "permission is required"}`)
			return
		}
		serviceAccountID := mux.Vars(r)["id"]
		has, err :=
			serviceAccountsUseCase.HasPermission(serviceAccountID, permissionSl[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !has {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
