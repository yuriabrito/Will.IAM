package usecases

import (
	"github.com/ghostec/Will.IAM/repositories"
)

// ServiceAccounts define entrypoints for ServiceAccount actions
type ServiceAccounts interface {
	Authenticate(string, string) error
	HasPermission(string) bool
}

type serviceAccounts struct {
	serviceAccountsRepository repositories.ServiceAccounts
	permissionsRepository     repositories.Permissions
	//	rolesRepository repositories.Roles
}

// NewServiceAccounts serviceAccounts ctor
func NewServiceAccounts(
	serviceAccountsRepository repositories.ServiceAccounts,
	permissionsRepository repositories.Permissions,
) ServiceAccounts {
	return &serviceAccounts{
		serviceAccountsRepository: serviceAccountsRepository,
		permissionsRepository:     permissionsRepository,
	}
}

// Authenticate verifies if token is valid for id, and sometimes refreshes it
func (u *serviceAccounts) Authenticate(id, token string) error {
	return nil
}

// HasPermission checks if user has the ownership level required to take an
// action over a resource
func (u serviceAccounts) HasPermission(permission string) bool {
	return false
}
