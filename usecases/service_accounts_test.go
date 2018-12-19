// +build integration

package usecases_test

import (
	"fmt"
	"testing"

	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/oauth2"
	"github.com/ghostec/Will.IAM/repositories"
	helpers "github.com/ghostec/Will.IAM/testing"
	"github.com/ghostec/Will.IAM/usecases"
)

func beforeEachServiceAccounts(t *testing.T) {
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

func TestServiceAccountsCreate(t *testing.T) {
	beforeEachServiceAccounts(t)
	saUC := getServiceAccountsUseCase(t)
	saM := &models.ServiceAccount{
		Name:  "some name",
		Email: "test@domain.com",
	}
	if err := saUC.Create(saM); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if saM.ID == "" {
		t.Errorf("Expected saM.ID to be non-empty")
	}
}

func TestServiceAccountsCreateShouldCreateRoleAndRoleBinding(t *testing.T) {
	beforeEachServiceAccounts(t)
	saUC := getServiceAccountsUseCase(t)
	saM := &models.ServiceAccount{
		Name:  "some name",
		Email: "test@domain.com",
	}
	if err := saUC.Create(saM); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	rs, err := saUC.GetRoles(saM.ID)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if len(rs) != 1 {
		t.Errorf("Should have only 1 role binding. Found %d", len(rs))
	}
	rName := fmt.Sprintf("service-account:%s", saM.ID)
	if rs[0].Name != rName {
		t.Errorf("Expected role name to be %s. Got %s", rName, rs[0].Name)
	}
}
