package api

import (
	"fmt"
	"net/http"

	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/usecases"
	"github.com/gorilla/mux"
)

func permissionsDeleteHandler(
	sasUC usecases.ServiceAccounts, psUC usecases.Permissions,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pID := mux.Vars(r)["id"]
		p, err := psUC.Get(pID)
		if err != nil {
			// TODO: use appropriate errors
			if err.Error() == fmt.Sprintf("permission %s not found", pID) {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		p.OwnershipLevel = models.OwnershipLevels.Owner
		saID, _ := getServiceAccountID(r.Context())
		has, err := sasUC.HasPermission(saID, p.ToString())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !has {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		err = psUC.Delete(pID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
