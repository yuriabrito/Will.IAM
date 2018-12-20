package usecases

import (
	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/repositories"
)

// Roles define entrypoints for ServiceAccount actions
type Roles interface {
	Create(r *models.Role) error
	CreatePermission(string, *models.Permission) error
	GetPermissions(string) ([]models.Permission, error)
}

type roles struct {
	rolesRepository       repositories.Roles
	permissionsRepository repositories.Permissions
}

func (rs roles) Create(r *models.Role) error {
	return rs.rolesRepository.Create(r)
}

func (rs roles) CreatePermission(roleID string, p *models.Permission) error {
	p.RoleID = roleID
	return rs.permissionsRepository.Create(p)
}

func (rs roles) GetPermissions(roleID string) ([]models.Permission, error) {
	r := models.Role{ID: roleID}
	return rs.permissionsRepository.ForRoles([]models.Role{r})
}

// NewRoles ctor
func NewRoles(
	rsRepo repositories.Roles, psRepo repositories.Permissions,
) Roles {
	return &roles{
		rolesRepository:       rsRepo,
		permissionsRepository: psRepo,
	}
}
