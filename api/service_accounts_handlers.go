package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/usecases"
	"github.com/gorilla/mux"
)

func serviceAccountsHasPermissionHandler(
	saUC usecases.ServiceAccounts,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		qs := r.URL.Query()
		permissionSl := qs["permission"]
		if len(permissionSl) == 0 {
			Write(w, http.StatusBadRequest, `{"error": "permission is required"}`)
			return
		}
		saID := mux.Vars(r)["id"]
		has, err :=
			saUC.HasPermission(saID, permissionSl[0])
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

func serviceAccountsGetHandler(
	saUC usecases.ServiceAccounts,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// qs := r.URL.Query()
		// withRoles := qs["withRoles"]
		// if len(withRoles) == 1 {
		// 	return
		// }
		saID := mux.Vars(r)["id"]
		sa, err := saUC.Get(saID)
		if err != nil {
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
	saUC usecases.ServiceAccounts,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		saID, _ := getServiceAccountID(r.Context())
		has, err := saUC.HasPermission(
			saID, models.BuildWillIAMPermissionStr(
				models.OwnershipLevels.Lender, "CreateServiceAccount", "*",
			),
		)
		if err != nil {
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
		_, err = saUC.CreateKeyPairType(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
