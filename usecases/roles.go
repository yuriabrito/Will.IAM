package usecases

import (
	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/repositories"
)

// RoleUpdate is the required data to update a role
type RoleUpdate struct {
	ID                 string              `json:"-"`
	Name               string              `json:"name"`
	PermissionsStrings []string            `json:"permissions"`
	Permissions        []models.Permission `json:"-"`
	ServiceAccountsIDs []string            `json:"service_accounts_ids"`
}

// Roles define entrypoints for ServiceAccount actions
type Roles interface {
	Create(r *models.Role) error
	CreatePermission(string, *models.Permission) error
	Update(RoleUpdate) error
	Get(string) (*models.Role, error)
	GetPermissions(string) ([]models.Permission, error)
	GetServiceAccounts(string) ([]models.ServiceAccount, error)
	WithNamePrefix(string, int) ([]models.Role, error)
	List() ([]models.Role, error)
}

type roles struct {
	repo *repositories.All
}

func (rs roles) Create(r *models.Role) error {
	return rs.repo.Roles.Create(r)
}

func (rs roles) CreatePermission(roleID string, p *models.Permission) error {
	p.RoleID = roleID
	return rs.repo.Permissions.Create(p)
}

func (rs roles) Update(ru RoleUpdate) error {
	// TODO: use tx
	for i := range ru.Permissions {
		if err := rs.CreatePermission(ru.ID, &ru.Permissions[i]); err != nil {
			return err
		}
	}
	role := &models.Role{ID: ru.ID, Name: ru.Name}
	return rs.repo.Roles.Update(role)
}

func (rs roles) GetPermissions(roleID string) ([]models.Permission, error) {
	return rs.repo.Permissions.ForRole(roleID)
}

func (rs roles) GetServiceAccounts(
	roleID string,
) ([]models.ServiceAccount, error) {
	return rs.repo.Roles.GetServiceAccounts(roleID)
}

func (rs roles) WithNamePrefix(
	prefix string, maxResults int,
) ([]models.Role, error) {
	return rs.repo.Roles.WithNamePrefix(prefix, maxResults)
}

func (rs roles) List() ([]models.Role, error) {
	return rs.repo.Roles.List()
}

func (rs roles) Get(id string) (*models.Role, error) {
	return rs.repo.Roles.Get(id)
}

// NewRoles ctor
func NewRoles(repo *repositories.All) Roles {
	return &roles{repo: repo}
}
