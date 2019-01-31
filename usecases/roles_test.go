// +build integration

package usecases_test

import (
	"testing"

	"github.com/ghostec/Will.IAM/models"
	helpers "github.com/ghostec/Will.IAM/testing"
)

func beforeEachRoles(t *testing.T) {
	t.Helper()
	storage := helpers.GetStorage(t)
	_, err := storage.PG.DB.Exec("TRUNCATE roles, permissions CASCADE;")
	if err != nil {
		panic(err)
	}
}

func TestRolesCreatePermission(t *testing.T) {
	beforeEachRoles(t)
	rsUC := helpers.GetRolesUseCase(t)
	pStr := "Maestro::RL::CreateScheduler::some-game::*"
	p, err := models.BuildPermission(pStr)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	r := &models.Role{Name: "Test role name"}
	if err := rsUC.Create(r); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	p.RoleID = r.ID
	if err := rsUC.CreatePermission(r.ID, &p); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	ps, err := rsUC.GetPermissions(p.RoleID)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if len(ps) != 1 {
		t.Errorf("Expected len(permissions) to be 1. Got %d", len(ps))
	}
	if ps[0].String() != pStr {
		t.Errorf("Expected permission to be %s. Got %s", pStr, ps[0].String())
	}
}
