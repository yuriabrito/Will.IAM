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
}
