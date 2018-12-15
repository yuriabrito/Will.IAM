package usecases

import (
	"fmt"

	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/repositories"
)

// ServiceAccounts define entrypoints for ServiceAccount actions
type ServiceAccounts interface {
	Create(*models.ServiceAccount) error
	Authenticate(string, string) error
	HasPermission(string, string) (bool, error)
	GetPermissions(string) ([]models.Permission, error)
	GetRoles(string) ([]models.Role, error)
}

type serviceAccounts struct {
	serviceAccountsRepository repositories.ServiceAccounts
	rolesRepository           repositories.Roles
	permissionsRepository     repositories.Permissions
}

// NewServiceAccounts serviceAccounts ctor
func NewServiceAccounts(
	serviceAccountsRepository repositories.ServiceAccounts,
	rolesRepository repositories.Roles,
	permissionsRepository repositories.Permissions,
) ServiceAccounts {
	return &serviceAccounts{
		serviceAccountsRepository: serviceAccountsRepository,
		rolesRepository:           rolesRepository,
		permissionsRepository:     permissionsRepository,
	}
}

func (sas serviceAccounts) Create(sa *models.ServiceAccount) error {
	// TODO: pass tx to repo create -> service_accounts + roles + role_bindings
	if err := sas.serviceAccountsRepository.Create(sa); err != nil {
		return err
	}
	r := &models.Role{
		Name: fmt.Sprintf("service-account:%s", sa.ID),
	}
	if err := sas.rolesRepository.Create(r); err != nil {
		return err
	}
	if err := sas.rolesRepository.Bind(*r, *sa); err != nil {
		return err
	}
	return nil
}

// GetRoles returns all roles to which the serviceAccountID is bound to
func (sas serviceAccounts) GetRoles(
	serviceAccountID string,
) ([]models.Role, error) {
	roles, err := sas.rolesRepository.ForServiceAccountID(serviceAccountID)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// Authenticate verifies if token is valid for id, and sometimes refreshes it
func (sas *serviceAccounts) Authenticate(id, token string) error {
	return nil
}

// HasPermission checks if user has the ownership level required to take an
// action over a resource
func (sas serviceAccounts) HasPermission(
	serviceAccountID, permissionStr string,
) (bool, error) {
	permissions, err := sas.GetPermissions(serviceAccountID)
	if err != nil {
		return false, err
	}
	permission, err := models.BuildPermission(permissionStr)
	if err != nil {
		return false, err
	}
	return permission.IsPresent(permissions), nil
}

func (sas serviceAccounts) GetPermissions(
	serviceAccountID string,
) ([]models.Permission, error) {
	roles, err := sas.rolesRepository.ForServiceAccountID(serviceAccountID)
	if err != nil {
		return nil, err
	}
	permissions, err := sas.permissionsRepository.ForRoles(roles)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
