package usecases

import (
	"fmt"

	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/oauth2"
	"github.com/ghostec/Will.IAM/repositories"
	"github.com/gofrs/uuid"
)

// ServiceAccounts define entrypoints for ServiceAccount actions
type ServiceAccounts interface {
	Create(*models.ServiceAccount) error
	AuthenticateAccessToken(string) (*AccessTokenAuth, error)
	AuthenticateKeyPair(string, string) (string, error)
	HasPermission(string, string) (bool, error)
	GetPermissions(string) ([]models.Permission, error)
	CreatePermission(string, models.Permission) error
	GetRoles(string) ([]models.Role, error)
}

type serviceAccounts struct {
	serviceAccountsRepository repositories.ServiceAccounts
	rolesRepository           repositories.Roles
	permissionsRepository     repositories.Permissions
	oauth2Provider            oauth2.Provider
}

// NewServiceAccounts serviceAccounts ctor
func NewServiceAccounts(
	serviceAccountsRepository repositories.ServiceAccounts,
	rolesRepository repositories.Roles,
	permissionsRepository repositories.Permissions,
	provider oauth2.Provider,
) ServiceAccounts {
	return &serviceAccounts{
		serviceAccountsRepository: serviceAccountsRepository,
		rolesRepository:           rolesRepository,
		permissionsRepository:     permissionsRepository,
		oauth2Provider:            provider,
	}
}

func (sas serviceAccounts) Create(sa *models.ServiceAccount) error {
	// TODO: pass tx to repo create -> service_accounts + roles + role_bindings
	sa.ID = uuid.Must(uuid.NewV4()).String()
	r := &models.Role{
		Name: fmt.Sprintf("service-account:%s", sa.ID),
	}
	if err := sas.rolesRepository.Create(r); err != nil {
		return err
	}
	sa.BaseRoleID = r.ID
	if err := sas.serviceAccountsRepository.Create(sa); err != nil {
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

// AccessTokenAuth stores a ServiceAccountID and the (maybe refreshed)
// AccessToken
type AccessTokenAuth struct {
	ServiceAccountID string
	AccessToken      string
}

// AuthenticateAccessToken verifies if token is valid for email, and sometimes refreshes it
func (sas *serviceAccounts) AuthenticateAccessToken(
	accessToken string,
) (*AccessTokenAuth, error) {
	authResult, err := sas.oauth2Provider.Authenticate(accessToken)
	if err != nil {
		return nil, err
	}
	sa, err := sas.serviceAccountsRepository.ForEmail(authResult.Email)
	if err != nil {
		return nil, err
	}
	return &AccessTokenAuth{
		ServiceAccountID: sa.ID,
		AccessToken:      authResult.AccessToken,
	}, nil
}

// AuthenticateKeyPair verifies if key pair is valid
func (sas *serviceAccounts) AuthenticateKeyPair(
	keyID, keySecret string,
) (string, error) {
	sa, err := sas.serviceAccountsRepository.ForKeyPair(keyID, keySecret)
	if err != nil {
		return "", err
	}
	return sa.ID, nil
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

func (sas serviceAccounts) CreatePermission(
	serviceAccountID string, permission models.Permission,
) error {
	sa, err := sas.serviceAccountsRepository.Get(serviceAccountID)
	if err != nil {
		return err
	}
	permission.RoleID = sa.BaseRoleID
	if err := sas.permissionsRepository.Create(permission); err != nil {
		return err
	}
	return nil
}
