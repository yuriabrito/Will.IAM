package usecases

import (
	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/repositories"
)

// Permissions define entrypoints for Permissions actions
type Permissions interface {
	Get(string) (*models.Permission, error)
	Delete(string) error
	Create(*models.Permission) error
	CreateRequest(string, *models.PermissionRequest) error
	GetPermissionRequests(string) ([]models.PermissionRequest, error)
	ForRoles([]models.Role) ([]models.Permission, error)
}

type permissions struct {
	permissionsRepository repositories.Permissions
}

func (ps permissions) Get(id string) (*models.Permission, error) {
	return ps.permissionsRepository.Get(id)
}

func (ps permissions) Delete(id string) error {
	return ps.permissionsRepository.Delete(id)
}

func (ps permissions) CreateRequest(
	saID string, r *models.PermissionRequest,
) error {
	r.State = models.PermissionRequestStates.Created
	return ps.permissionsRepository.CreateRequest(saID, r)
}

func (ps permissions) Create(p *models.Permission) error {
	return ps.permissionsRepository.Create(p)
}

func (ps permissions) ForRoles(
	rs []models.Role,
) ([]models.Permission, error) {
	return ps.permissionsRepository.ForRoles(rs)
}

func (ps permissions) GetPermissionRequests(
	saID string,
) ([]models.PermissionRequest, error) {
	return ps.permissionsRepository.GetPermissionRequests(saID)
}

// NewPermissions ctor
func NewPermissions(psRepo repositories.Permissions) Permissions {
	return &permissions{
		permissionsRepository: psRepo,
	}
}
