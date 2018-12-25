// +build integration

package api_test

import (
	"fmt"
	"net/http"
	"testing"

	helpers "github.com/ghostec/Will.IAM/testing"
	"github.com/gofrs/uuid"
)

func beforeEachPermissionsHandlers(t *testing.T) {
	t.Helper()
	storage := helpers.GetStorage(t)
	rels := []string{"permissions", "role_bindings", "service_accounts", "roles"}
	for _, rel := range rels {
		if _, err := storage.PG.DB.Exec(
			fmt.Sprintf("DELETE FROM %s", rel),
		); err != nil {
			panic(err)
		}
	}
}

func TestPermissionsDeleteHandlerNonExistentID(t *testing.T) {
	beforeEachRolesHandlers(t)
	rootSA := helpers.CreateRootServiceAccount(t)
	app := helpers.GetApp(t)
	req, _ := http.NewRequest("DELETE", fmt.Sprintf(
		"/permissions/%s", uuid.Must(uuid.NewV4()).String(),
	), nil)
	req.Header.Set("Authorization", fmt.Sprintf(
		"KeyPair %s:%s", rootSA.KeyID, rootSA.KeySecret,
	))
	rec := helpers.DoRequest(t, req, app.GetRouter())
	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status 204. Got %d", rec.Code)
	}
}

func TestPermissionsDeleteHandlerNonRootSA(t *testing.T) {
	beforeEachRolesHandlers(t)
	saUC := helpers.GetServiceAccountsUseCase(t)
	rootSA := helpers.CreateRootServiceAccount(t)
	sa, err := saUC.CreateKeyPairType("some sa")
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
		return
	}
	app := helpers.GetApp(t)
	p := "SomeService::RO::SomeAction::*"
	req, _ := http.NewRequest("POST", fmt.Sprintf(
		"/roles/%s/permissions?permission=%s", sa.BaseRoleID, p,
	), nil)
	req.Header.Set("Authorization", fmt.Sprintf(
		"KeyPair %s:%s", rootSA.KeyID, rootSA.KeySecret,
	))
	rec := helpers.DoRequest(t, req, app.GetRouter())
	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status 201. Got %d", rec.Code)
	}
	permissions, err := saUC.GetPermissions(sa.ID)
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
		return
	}
	if len(permissions) != 1 {
		t.Errorf("Expected only 1 permission. Got %d", len(permissions))
		return
	}
	deleterSA, err := saUC.CreateKeyPairType("deleter sa")
	req, _ = http.NewRequest("DELETE", fmt.Sprintf(
		"/permissions/%s", permissions[0].ID,
	), nil)
	req.Header.Set("Authorization", fmt.Sprintf(
		"KeyPair %s:%s", deleterSA.KeyID, deleterSA.KeySecret,
	))
	rec = helpers.DoRequest(t, req, app.GetRouter())
	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status 403. Got %d", rec.Code)
	}
}

func TestPermissionsDeleteHandler(t *testing.T) {
	beforeEachRolesHandlers(t)
	saUC := helpers.GetServiceAccountsUseCase(t)
	rootSA := helpers.CreateRootServiceAccount(t)
	sa, err := saUC.CreateKeyPairType("some sa")
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
		return
	}
	app := helpers.GetApp(t)
	p := "SomeService::RO::SomeAction::*"
	req, _ := http.NewRequest("POST", fmt.Sprintf(
		"/roles/%s/permissions?permission=%s", sa.BaseRoleID, p,
	), nil)
	req.Header.Set("Authorization", fmt.Sprintf(
		"KeyPair %s:%s", rootSA.KeyID, rootSA.KeySecret,
	))
	rec := helpers.DoRequest(t, req, app.GetRouter())
	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status 201. Got %d", rec.Code)
	}
	permissions, err := saUC.GetPermissions(sa.ID)
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
		return
	}
	if len(permissions) != 1 {
		t.Errorf("Expected only 1 permission. Got %d", len(permissions))
		return
	}
	req, _ = http.NewRequest("DELETE", fmt.Sprintf(
		"/permissions/%s", permissions[0].ID,
	), nil)
	req.Header.Set("Authorization", fmt.Sprintf(
		"KeyPair %s:%s", rootSA.KeyID, rootSA.KeySecret,
	))
	rec = helpers.DoRequest(t, req, app.GetRouter())
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200. Got %d", rec.Code)
	}
}