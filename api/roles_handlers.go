package api

import (
	"net/http"

	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/usecases"
	"github.com/gorilla/mux"
)

func rolesCreatePermissionHandler(
	sasUC usecases.ServiceAccounts, rsUC usecases.Roles,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		qs := r.URL.Query()
		permissionSl := qs["permission"]
		if len(permissionSl) == 0 {
			Write(w, http.StatusUnprocessableEntity, `{"error": "querystrings.permission is required"}`)
			return
		}

		saID, _ := getServiceAccountID(r.Context())
		has, err := sasUC.HasPermission(saID, permissionSl[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !has {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		rID := mux.Vars(r)["id"]
		p, err := models.BuildPermission(permissionSl[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = rsUC.CreatePermission(rID, &p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
