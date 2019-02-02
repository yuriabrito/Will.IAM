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
		p, err := psUC.WithCtx(r.Context()).Get(pID)
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
		has, err := sasUC.WithCtx(r.Context()).HasPermissionString(saID, p.String())
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !has {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		err = psUC.WithCtx(r.Context()).Delete(pID)
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
		has, err := sasUC.WithCtx(r.Context()).
			HasPermissionString(saID, pr.ToLenderString())
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

		err = psUC.WithCtx(r.Context()).CreateRequest(saID, pr)
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
		prs, err := psUC.WithCtx(r.Context()).GetPermissionRequests(saID)
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

func permissionsHasHandler(
	sasUC usecases.ServiceAccounts,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
		qs := r.URL.Query()
		permissionSl := qs["permission"]
		if len(permissionSl) == 0 {
			Write(w, http.StatusUnprocessableEntity,
				`{"error": "querystrings.permission is required"}`)
			return
		}
		saID, _ := getServiceAccountID(r.Context())
		has, err :=
			sasUC.WithCtx(r.Context()).HasPermissionString(saID, permissionSl[0])
		if err != nil {
			l.Error(err)
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
