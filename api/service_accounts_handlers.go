package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/usecases"
	"github.com/gorilla/mux"
	"github.com/topfreegames/extensions/middleware"
)

func serviceAccountsGetHandler(
	sasUC usecases.ServiceAccounts,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
		// qs := r.URL.Query()
		// withRoles := qs["withRoles"]
		// if len(withRoles) == 1 {
		// 	return
		// }
		saID := mux.Vars(r)["id"]
		sa, err := sasUC.Get(saID)
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data := map[string]string{
			"id":         sa.ID,
			"name":       sa.Name,
			"baseRoleId": sa.BaseRoleID,
		}
		bts, err := json.Marshal(data)
		WriteBytes(w, 200, bts)
	}
}

func serviceAccountsCreateHandler(
	sasUC usecases.ServiceAccounts,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
		saID, _ := getServiceAccountID(r.Context())
		has, err := sasUC.HasPermission(
			saID, models.BuildWillIAMPermissionStr(
				models.OwnershipLevels.Lender, "CreateServiceAccount", "*",
			),
		)
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !has {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		m := map[string]interface{}{}
		err = json.Unmarshal(body, &m)
		name, ok := m["name"].(string)
		if !ok || name == "" {
			Write(w, http.StatusUnprocessableEntity, "body.name is required")
			return
		}
		_, err = sasUC.CreateKeyPairType(name)
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func serviceAccountsListHandler(
	sasUC usecases.ServiceAccounts,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
		sa, err := sasUC.Get("asdf")
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data := map[string]string{
			"id":         sa.ID,
			"name":       sa.Name,
			"baseRoleId": sa.BaseRoleID,
		}
		bts, err := json.Marshal(data)
		WriteBytes(w, 200, bts)
	}
}
