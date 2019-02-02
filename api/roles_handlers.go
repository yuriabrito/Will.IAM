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
		has, err := sasUC.WithCtx(r.Context()).
			HasPermissionString(saID, sameP.String())
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
		err = rsUC.WithCtx(r.Context()).CreatePermission(rID, &p)
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
		ru := usecases.RoleUpdate{}
		err = json.Unmarshal(body, &ru)
		if err != nil {
			l.WithError(err).Error("rolesUpdateHandler json.Unmarshal(body)")
			Write(w, http.StatusBadRequest, `{"error": "body malformed"}`)
			return
		}
		saID, _ := getServiceAccountID(r.Context())
		ru.Permissions, err = models.BuildPermissions(ru.PermissionsStrings)
		if err != nil {
			l.WithError(err).Error("rolesUpdateHandler models.BuildPermissions")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		has, err := sasUC.WithCtx(r.Context()).
			HasAllOwnerPermissions(saID, ru.Permissions)
		if err != nil {
			l.WithError(err).Error("rolesUpdateHandler sasUC.HasAllOwnerPermissionsStrings")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !has {
			Write(
				w, http.StatusForbidden,
				`{ "error": "not owner of all permissions" }`,
			)
			return
		}
		ru.ID = mux.Vars(r)["id"]
		if err = rsUC.WithCtx(r.Context()).Update(ru); err != nil {
			l.WithError(err).Error("rolesUpdateHandler rsUC.Update")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// TODO: audit
		// TODO: GetPermissions and delete diff
		w.WriteHeader(http.StatusOK)
	}
}

func rolesListHandler(
	rsUC usecases.Roles,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
		rsSl, err := rsUC.WithCtx(r.Context()).List()
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
		err = rsUC.WithCtx(r.Context()).Create(role)
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
		rsUCc := rsUC.WithCtx(r.Context())
		role, err := rsUCc.Get(id)
		if err != nil {
			l.WithError(err).Error("rolesViewHandler rsUC.Get")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		pSl, err := rsUCc.GetPermissions(id)
		if err != nil {
			l.WithError(err).Error("rolesViewHandler rsUC.GetPermissions")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		permissions := make([]string, len(pSl))
		for i := range pSl {
			permissions[i] = pSl[i].String()
		}
		sas, err := rsUCc.GetServiceAccounts(id)
		if err != nil {
			l.WithError(err).Error("rolesViewHandler rsUC.GetServiceAccounts")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		sasFiltered, err := keepJSONFields(sas, "id", "name", "picture", "email")
		if err != nil {
			l.WithError(err).Error("rolesViewHandler keepJSONFields(sas)")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		body := map[string]interface{}{
			"id":              role.ID,
			"name":            role.Name,
			"permissions":     permissions,
			"serviceAccounts": sasFiltered,
		}
		bts, err := json.Marshal(body)
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		WriteBytes(w, 200, bts)
	}
}
