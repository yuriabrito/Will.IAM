// +build integration

package api_test

import (
	"net/http"
	"testing"

	helpers "github.com/ghostec/Will.IAM/testing"
)

func TestHasPermission(t *testing.T) {
	type hasPermissionTest struct {
		request        string
		expectedStatus int
	}
	tt := []hasPermissionTest{
		hasPermissionTest{
			request:        "/service_accounts/1234/permissions",
			expectedStatus: http.StatusBadRequest,
		},
		hasPermissionTest{
			request:        "/service_accounts/1234/permissions?permission=Service::RL::TestAction::*",
			expectedStatus: http.StatusForbidden,
		},
	}

	app := helpers.GetApp(t)
	for _, tt := range tt {
		req, _ := http.NewRequest("GET", tt.request, nil)
		rec := helpers.DoRequest(t, req, app.GetRouter())
		if rec.Code != tt.expectedStatus {
			t.Errorf("Expected http.StatusForbidden. Got %d", rec.Code)
		}
	}
}
