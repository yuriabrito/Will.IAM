// +build integration

package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/oauth2"
	"github.com/ghostec/Will.IAM/repositories"
	helpers "github.com/ghostec/Will.IAM/testing"
	"github.com/ghostec/Will.IAM/usecases"
)

func beforeEachServices(t *testing.T) {
	t.Helper()
	storage := helpers.GetStorage(t)
	_, err := storage.PG.DB.Exec("TRUNCATE service_accounts CASCADE;")
	if err != nil {
		panic(err)
	}
}

func getServiceAccountsUseCase(t *testing.T) usecases.ServiceAccounts {
	t.Helper()
	storage := helpers.GetStorage(t)
	saRepo := repositories.NewServiceAccounts(storage)
	rRepo := repositories.NewRoles(storage)
	pRepo := repositories.NewPermissions(storage)
	providerBlankMock := oauth2.NewProviderBlankMock()
	return usecases.NewServiceAccounts(saRepo, rRepo, pRepo, providerBlankMock)
}

func TestServicesCreateHandler(t *testing.T) {
	beforeEachServices(t)
	saUC := getServiceAccountsUseCase(t)
	sa := &models.ServiceAccount{
		Name:  "any",
		Email: "any@email.com",
	}
	if err := saUC.Create(sa); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
		return
	}
	type createTest struct {
		service        *models.Service
		expectedStatus int
	}
	tt := []createTest{
		createTest{
			service: &models.Service{
				Name:                    "Some Service",
				PermissionName:          "SomeService",
				CreatorServiceAccountID: sa.ID,
			},
			expectedStatus: http.StatusCreated,
		},
	}

	app := helpers.GetApp(t)
	for _, tt := range tt {
		bts, err := json.Marshal(tt.service)
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
			return
		}
		req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(bts))
		req.Header.Set("Authorization", "Bearer dummy_access_token")
		rec := helpers.DoRequest(t, req, app.GetRouter())
		if rec.Code != tt.expectedStatus {
			t.Errorf("Expected %d. Got %d", tt.expectedStatus, rec.Code)
		}
	}
}
