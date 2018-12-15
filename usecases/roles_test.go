// +build integration

package usecases_test

import (
	"testing"

	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/repositories"
	helpers "github.com/ghostec/Will.IAM/testing"
	"github.com/ghostec/Will.IAM/usecases"
)

func beforeEachRoles(t *testing.T) {
	t.Helper()
	storage := helpers.GetStorage(t)
	_, err := storage.PG.DB.Exec("TRUNCATE roles, permissions CASCADE;")
	if err != nil {
		panic(err)
	}
}

func getRolesUseCase(t *testing.T) usecases.Roles {
	t.Helper()
	storage := helpers.GetStorage(t)
	rsRepo := repositories.NewRoles(storage)
	psRepo := repositories.NewPermissions(storage)
	return usecases.NewRoles(rsRepo, psRepo)
}

func TestRolesCreatePermission(t *testing.T) {
	beforeEachRoles(t)
	rsUC := getRolesUseCase(t)
	pStr := "Maestro::RL::CreateScheduler::some-game::*"
	p, err := models.BuildPermission(pStr)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	p.RoleID = "798a4f9a-cad4-4bd5-86a6-ca6a99fa73d5"
	if err := rsUC.CreatePermission(p); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	ps, err := rsUC.GetPermissions(p.RoleID)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if len(ps) != 1 {
		t.Errorf("Expected len(permissions) to be 1. Got %d", len(ps))
	}
	if ps[0].ToString() != pStr {
		t.Errorf("Expected permission to be %s. Got %s", pStr, ps[0].ToString())
	}
}
