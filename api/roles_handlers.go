package api

import (
	"net/http"

	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/usecases"
	"github.com/gorilla/mux"
	"github.com/topfreegames/extensions/middleware"
)

func rolesCreatePermissionHandler(
	sasUC usecases.ServiceAccounts, rsUC usecases.Roles,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
		qs := r.URL.Query()
		permissionSl := qs["permission"]
		if len(permissionSl) == 0 {
			Write(w, http.StatusUnprocessableEntity, `{"error": "querystrings.permission is required"}`)
			return
		}
		sameP, err := models.BuildPermission(permissionSl[0])
		if err != nil {
			Write(w, http.StatusBadRequest, `{"error": "querystrings.permission malformed"}`)
			return
		}
		sameP.OwnershipLevel = models.OwnershipLevels.Owner

		saID, _ := getServiceAccountID(r.Context())
		has, err := sasUC.HasPermission(saID, sameP.ToString())
		if err != nil {
			l.Error(err)
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
			Write(w, http.StatusBadRequest, `{"error": "querystrings.permission malformed"}`)
			return
		}
		err = rsUC.CreatePermission(rID, &p)
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
