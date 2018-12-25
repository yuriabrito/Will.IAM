// +build integration

package api_test

import (
	"fmt"
	"net/http"
	"testing"

	helpers "github.com/ghostec/Will.IAM/testing"
)

func beforeEachRolesHandlers(t *testing.T) {
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

func TestRolesCreatePermissionHandler(t *testing.T) {
	beforeEachRolesHandlers(t)
	rootSA := helpers.CreateRootServiceAccount(t)
	saUC := helpers.GetServiceAccountsUseCase(t)
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
	if permissions[0].ToString() != p {
		t.Errorf(
			"Expected permission to be %s. Got %s", p, permissions[0].ToString(),
		)
		return
	}
}

func TestRolesCreatePermissionHandlerNonRootSA(t *testing.T) {
	beforeEachRolesHandlers(t)
	saUC := helpers.GetServiceAccountsUseCase(t)
	creatorSA, err := saUC.CreateKeyPairType("creator sa")
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
		return
	}
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
		"KeyPair %s:%s", creatorSA.KeyID, creatorSA.KeySecret,
	))
	rec := helpers.DoRequest(t, req, app.GetRouter())
	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status 403. Got %d", rec.Code)
	}
}
