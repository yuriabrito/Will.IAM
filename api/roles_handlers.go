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
		has, err := sasUC.HasPermission(saID, sameP.String())
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

func rolesUpdateHandler(
	sasUC usecases.ServiceAccounts, rsUC usecases.Roles,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			l.WithError(err).Error("rolesUpdateHandler ioutil.ReadAll(body)")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		m := map[string]interface{}{}
		err = json.Unmarshal(body, &m)
		if err != nil {
			l.WithError(err).Error("rolesUpdateHandler json.Unmarshal(body)")
			Write(w, http.StatusBadRequest, `{"error": "body malformed"}`)
			return
		}
		roleID := mux.Vars(r)["id"]
		name, ok := m["name"]
		if !ok {
			l.WithError(err).Error("rolesUpdateHandler name is blank")
			Write(w, http.StatusUnprocessableEntity, `{"error": "name is required"}`)
			return
		}
		if _, ok = name.(string); !ok {
			l.WithError(err).Error("rolesUpdateHandler name must be a string")
			Write(
				w, http.StatusUnprocessableEntity, `{"error": "name must be a string"}`,
			)
			return
		}
		// TODO: use tx
		role := &models.Role{ID: roleID, Name: name.(string)}
		if err = rsUC.Update(role); err != nil {
			l.WithError(err).Error("rolesUpdateHandler rsUC.Update")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// TODO: audit
		permissionsI, ok := m["permissions"]
		if !ok {
			w.WriteHeader(http.StatusOK)
			return
		}
		permissions, ok := permissionsI.([]interface{})
		if !ok {
			l.WithError(err).Error(
				"rolesUpdateHandler permissions must be an array of strings",
			)
			Write(
				w, http.StatusUnprocessableEntity,
				`{"error": "permissions must be an array of strings"}`,
			)
			return
		}

		pSl := make([]models.Permission, len(permissions))
		for i := range permissions {
			pStr, ok := permissions[i].(string)
			if !ok {
				Write(
					w, http.StatusUnprocessableEntity, `{"error": "permission malformed"}`,
				)
				return
			}
			sameP, err := models.BuildPermission(pStr)
			if err != nil {
				Write(
					w, http.StatusUnprocessableEntity, `{"error": "permission malformed"}`,
				)
				return
			}
			sameP.OwnershipLevel = models.OwnershipLevels.Owner
			pSl[i] = sameP
		}

		saID, _ := getServiceAccountID(r.Context())
		has, err := sasUC.HasPermissions(saID, pSl)
		if err != nil {
			l.WithError(err).Error("rolesUpdateHandler sasUC.HasPermissions")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range has {
			if !has[i] {
				Write(
					w, http.StatusForbidden,
					fmt.Sprintf(
						`{ "error": "not owner of %s" }`,
						m["permissions"].([]interface{})[i].(string),
					),
				)
				return
			}
		}
		for i := range pSl {
			if err := rsUC.CreatePermission(roleID, &pSl[i]); err != nil {
				l.WithError(err).Error("rolesUpdateHandler rsUC.CreatePermission")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func rolesListHandler(
	rsUC usecases.Roles,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
		rsSl, err := rsUC.List()
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		bts, err := keepJSONFieldsBytes(rsSl, "id", "name")
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		WriteBytes(w, 200, bts)
	}
}

func rolesCreateHandler(
	rsUC usecases.Roles,
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
		m := map[string]interface{}{}
		err = json.Unmarshal(body, &m)
		name, ok := m["name"].(string)
		if !ok || name == "" {
			Write(w, http.StatusUnprocessableEntity,
				`{ "error": { "name": "required" } }`)
			return
		}
		role := &models.Role{Name: name, IsBaseRole: false}
		err = rsUC.Create(role)
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func rolesViewHandler(
	rsUC usecases.Roles,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
		id := mux.Vars(r)["id"]
		role, err := rsUC.Get(id)
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		bts, err := json.Marshal(role)
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		WriteBytes(w, 200, bts)
	}
}
