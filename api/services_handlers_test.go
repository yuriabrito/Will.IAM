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

func getServicesUseCase(t *testing.T) usecases.Services {
	t.Helper()
	storage := helpers.GetStorage(t)
	sRepo := repositories.NewServices(storage)
	saUC := getServiceAccountsUseCase(t)
	return usecases.NewServices(sRepo, saUC)
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

	service := &models.Service{
		Name:                    "Some Service",
		PermissionName:          "SomeService",
		CreatorServiceAccountID: sa.ID,
		AMURL:                   "http://localhost:3333/am",
	}

	app := helpers.GetApp(t)
	bts, err := json.Marshal(service)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
		return
	}
	req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(bts))
	req.Header.Set("Authorization", "Bearer dummy_access_token")
	rec := helpers.DoRequest(t, req, app.GetRouter())
	if rec.Code != http.StatusCreated {
		t.Errorf("Expected 201. Got %d", rec.Code)
		return
	}
	sUC := getServicesUseCase(t)
	ss, err := sUC.All()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
		return
	}
	if len(ss) != 1 {
		t.Errorf("Expected to have 1 service. Got %d", len(ss))
		return
	}
	if ss[0].Name != "Some Service" {
		t.Errorf("Expected service name to be Some Service. Got %s", ss[0].Name)
		return
	}
	if ss[0].PermissionName != "SomeService" {
		t.Errorf(
			"Expected service permission name to be SomeService. Got %s",
			ss[0].PermissionName,
		)
		return
	}
}
