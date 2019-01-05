package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/usecases"
	"github.com/gorilla/mux"
	"github.com/topfreegames/extensions/middleware"
)

func permissionsDeleteHandler(
	sasUC usecases.ServiceAccounts, psUC usecases.Permissions,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
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
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !has {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		err = psUC.Delete(pID)
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func permissionsCreatePermissionRequestHandler(
	sasUC usecases.ServiceAccounts, psUC usecases.Permissions,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		pr := &models.PermissionRequest{}
		err = json.Unmarshal(body, pr)
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		saID, _ := getServiceAccountID(r.Context())
		has, err := sasUC.HasPermission(saID, pr.ToLenderString())
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if has {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// TODO: check if there's a request with state = Created already NoContent

		err = psUC.CreateRequest(saID, pr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func permissionsGetPermissionRequestsHandler(
	sasUC usecases.ServiceAccounts, psUC usecases.Permissions,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
		saID, _ := getServiceAccountID(r.Context())
		prs, err := psUC.GetPermissionRequests(saID)
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		bts, err := json.Marshal(prs)
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		WriteBytes(w, http.StatusOK, bts)
	}
}
