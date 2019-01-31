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
	CreateKeyPairType(string) (*models.ServiceAccount, error)
	CreateOAuth2Type(string, string) (*models.ServiceAccount, error)
	AuthenticateAccessToken(string) (*AccessTokenAuth, error)
	AuthenticateKeyPair(string, string) (string, error)
	HasPermissionString(string, string) (bool, error)
	HasAllOwnerPermissions(string, []models.Permission) (bool, error)
	HasPermissionsStrings(string, []string) ([]bool, error)
	HasPermissions(string, []models.Permission) ([]bool, error)
	GetPermissions(string) ([]models.Permission, error)
	CreatePermission(string, *models.Permission) error
	Get(string) (*models.ServiceAccount, error)
	List() ([]models.ServiceAccount, error)
	Search(string) ([]models.ServiceAccount, error)
	GetRoles(string) ([]models.Role, error)
}

type serviceAccounts struct {
	repo           *repositories.All
	oauth2Provider oauth2.Provider
}

// NewServiceAccounts serviceAccounts ctor
func NewServiceAccounts(
	repo *repositories.All,
	provider oauth2.Provider,
) ServiceAccounts {
	return &serviceAccounts{
		repo:           repo,
		oauth2Provider: provider,
	}
}

func (sas serviceAccounts) Create(sa *models.ServiceAccount) error {
	// TODO: pass tx to repo create -> service_accounts + roles + role_bindings
	sa.ID = uuid.Must(uuid.NewV4()).String()
	r := &models.Role{
		Name:       fmt.Sprintf("service-account:%s", sa.ID),
		IsBaseRole: true,
	}
	if err := sas.repo.Roles.Create(r); err != nil {
		return err
	}
	sa.BaseRoleID = r.ID
	if err := sas.repo.ServiceAccounts.Create(sa); err != nil {
		return err
	}
	if err := sas.repo.Roles.Bind(*r, *sa); err != nil {
		return err
	}
	return nil
}

// CreateKeyPairType will build a random key pair and create a
// Service Account with it
func (sas serviceAccounts) CreateKeyPairType(
	saName string,
) (*models.ServiceAccount, error) {
	// TODO: pass tx to repo create -> service_accounts + roles + role_bindings
	saKP := models.BuildKeyPairServiceAccount(saName)
	if err := sas.Create(saKP); err != nil {
		return nil, err
	}
	return saKP, nil
}

// CreateOAuth2Type creates an oauth2 service account
func (sas serviceAccounts) CreateOAuth2Type(
	saName, saEmail string,
) (*models.ServiceAccount, error) {
	saOAuth2 := models.BuildOAuth2ServiceAccount(saName, saEmail)
	if err := sas.Create(saOAuth2); err != nil {
		return nil, err
	}
	return saOAuth2, nil
}

// GetRoles returns all roles to which the serviceAccountID is bound to
func (sas serviceAccounts) GetRoles(
	serviceAccountID string,
) ([]models.Role, error) {
	roles, err := sas.repo.Roles.ForServiceAccountID(serviceAccountID)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// Get returns a service account by id
func (sas serviceAccounts) Get(
	serviceAccountID string,
) (*models.ServiceAccount, error) {
	sa, err := sas.repo.ServiceAccounts.Get(serviceAccountID)
	if err != nil {
		return nil, err
	}
	return sa, nil
}

// List returns a list of all service accounts
func (sas serviceAccounts) List() ([]models.ServiceAccount, error) {
	saSl, err := sas.repo.ServiceAccounts.List()
	if err != nil {
		return nil, err
	}
	return saSl, nil
}

// Search searches over Service Accounts names and emails
func (sas serviceAccounts) Search(
	term string,
) ([]models.ServiceAccount, error) {
	saSl, err := sas.repo.ServiceAccounts.Search(term)
	if err != nil {
		return nil, err
	}
	return saSl, nil
}

// AccessTokenAuth stores a ServiceAccountID and the (maybe refreshed)
// AccessToken
type AccessTokenAuth struct {
	ServiceAccountID string
	AccessToken      string
	Email            string
}

// AuthenticateAccessToken verifies if token is valid for email, and sometimes refreshes it
func (sas *serviceAccounts) AuthenticateAccessToken(
	accessToken string,
) (*AccessTokenAuth, error) {
	authResult, err := sas.oauth2Provider.Authenticate(accessToken)
	if err != nil {
		return nil, err
	}
	if authResult.Picture != "" {
		sa := &models.ServiceAccount{
			Name:    authResult.Email,
			Email:   authResult.Email,
			Picture: authResult.Picture,
		}
		if err = sas.Create(sa); err != nil {
			return nil, err
		}
	}
	sa, err := sas.repo.ServiceAccounts.ForEmail(authResult.Email)
	// TODO: use better errors
	saNotFoundErr := fmt.Sprintf(
		"Service Account not found for email %s", authResult.Email,
	)
	if err != nil && err.Error() != saNotFoundErr {
		return nil, err
	}
	if err != nil && err.Error() == saNotFoundErr {
	}
	return &AccessTokenAuth{
		ServiceAccountID: sa.ID,
		AccessToken:      authResult.AccessToken,
		Email:            authResult.Email,
	}, nil
}

// AuthenticateKeyPair verifies if key pair is valid
func (sas *serviceAccounts) AuthenticateKeyPair(
	keyID, keySecret string,
) (string, error) {
	sa, err := sas.repo.ServiceAccounts.ForKeyPair(keyID, keySecret)
	if err != nil {
		return "", err
	}
	return sa.ID, nil
}

// HasPermissionString checks if user has the ownership level required to take an
// action over a resource
func (sas serviceAccounts) HasPermissionString(
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

func (sas serviceAccounts) HasAllOwnerPermissions(
	serviceAccountID string, permissions []models.Permission,
) (bool, error) {
	for i := range permissions {
		permissions[i].OwnershipLevel = models.OwnershipLevels.Owner
	}
	has, err := sas.HasPermissions(serviceAccountID, permissions)
	if err != nil {
		return false, err
	}
	for i := range has {
		if !has[i] {
			return false, nil
		}
	}
	return true, nil
}

// HasPermissionsStrings returns an array of bools indicating whether a service
// account has some permissions
func (sas serviceAccounts) HasPermissionsStrings(
	serviceAccountID string, permissions []string,
) ([]bool, error) {
	pSl := make([]models.Permission, len(permissions))
	var err error
	for i := range permissions {
		pSl[i], err = models.BuildPermission(permissions[i])
		if err != nil {
			return nil, err
		}
	}
	return sas.HasPermissions(serviceAccountID, pSl)
}

// HasPermissions returns an array of bools indicating whether a service
// account has some permissions
func (sas serviceAccounts) HasPermissions(
	serviceAccountID string, permissions []models.Permission,
) ([]bool, error) {
	saPermissions, err := sas.GetPermissions(serviceAccountID)
	if err != nil {
		return nil, err
	}
	has := make([]bool, len(permissions))
	for i := range permissions {
		has[i] = permissions[i].IsPresent(saPermissions)
	}
	return has, nil
}

func (sas serviceAccounts) GetPermissions(
	serviceAccountID string,
) ([]models.Permission, error) {
	return sas.repo.Permissions.ForServiceAccount(serviceAccountID)
}

func (sas serviceAccounts) CreatePermission(
	serviceAccountID string, permission *models.Permission,
) error {
	sa, err := sas.repo.ServiceAccounts.Get(serviceAccountID)
	if err != nil {
		return err
	}
	permission.RoleID = sa.BaseRoleID
	if err := sas.repo.Permissions.Create(permission); err != nil {
		return err
	}
	return nil
}
