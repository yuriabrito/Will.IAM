// +build integration

package usecases_test

import (
	"testing"

	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/repositories"
	helpers "github.com/ghostec/Will.IAM/testing"
	"github.com/ghostec/Will.IAM/usecases"
)

func beforeEachServiceAccounts(t *testing.T) {
	t.Helper()
	storage := helpers.GetStorage(t)
	_, err := storage.PG.DB.Exec("DELETE FROM service_accounts;")
	if err != nil {
		panic(err)
	}
}

func getServiceAccountsUseCase(t *testing.T) usecases.ServiceAccounts {
	t.Helper()
	storage := helpers.GetStorage(t)
	saRepo := repositories.NewServiceAccounts(storage)
	pRepo := repositories.NewPermissions(storage)
	return usecases.NewServiceAccounts(saRepo, pRepo)
}

func TestServiceAccountsCreate(t *testing.T) {
	beforeEachServiceAccounts(t)
	saUC := getServiceAccountsUseCase(t)
	saM := &models.ServiceAccount{
		Email: "test@domain.com",
	}
	if err := saUC.Create(saM); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if saM.ID == "" {
		t.Errorf("Expected saM.ID to be non-empty")
	}
}
