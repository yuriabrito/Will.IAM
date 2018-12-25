// +build integration

package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	helpers "github.com/ghostec/Will.IAM/testing"
)

func beforeEachServiceAccountsHandlers(t *testing.T) {
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

func TestServiceAccountHasPermissionHandler(t *testing.T) {
	beforeEachServiceAccountsHandlers(t)
	type hasPermissionTest struct {
		request        string
		expectedStatus int
	}
	tt := []hasPermissionTest{
		hasPermissionTest{
			request:        "/service_accounts/1234/permissions",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		hasPermissionTest{
			request:        "/service_accounts/c91cfe41-a9cb-4685-8517-2411d1445616/permissions?permission=Service::RL::TestAction::*",
			expectedStatus: http.StatusForbidden,
		},
	}

	app := helpers.GetApp(t)
	for _, tt := range tt {
		req, _ := http.NewRequest("GET", tt.request, nil)
		rec := helpers.DoRequest(t, req, app.GetRouter())
		if rec.Code != tt.expectedStatus {
			t.Errorf("Expected %d. Got %d", tt.expectedStatus, rec.Code)
		}
	}
}

func TestServiceAccountCreateHandler(t *testing.T) {
	beforeEachServiceAccountsHandlers(t)
	type createTest struct {
		body           map[string]interface{}
		expectedStatus int
	}
	tt := []createTest{
		createTest{
			body: map[string]interface{}{
				"name": "some name",
			},
			expectedStatus: http.StatusCreated,
		},
		createTest{
			body:           map[string]interface{}{},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}
	rootSA := helpers.CreateRootServiceAccount(t)

	app := helpers.GetApp(t)
	for _, tt := range tt {
		bts, err := json.Marshal(tt.body)
		if err != nil {
			t.Errorf("Unexpected error %s", err.Error())
			return
		}
		req, _ := http.NewRequest("POST", "/service_accounts", bytes.NewBuffer(bts))
		req.Header.Set("Authorization", fmt.Sprintf(
			"KeyPair %s:%s", rootSA.KeyID, rootSA.KeySecret,
		))
		rec := helpers.DoRequest(t, req, app.GetRouter())
		if rec.Code != tt.expectedStatus {
			t.Errorf("Expected status %d. Got %d", tt.expectedStatus, rec.Code)
		}
	}
}
