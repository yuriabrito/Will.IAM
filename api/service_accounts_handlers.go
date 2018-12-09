package api

import (
	"net/http"

	"github.com/ghostec/Will.IAM/usecases"
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
		if has := serviceAccountsUseCase.HasPermission(permissionSl[0]); has {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusForbidden)
	}
}
