package usecases

import (
	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/repositories"
)

// Permissions define entrypoints for Permissions actions
type Permissions interface {
	Get(string) (*models.Permission, error)
	Delete(string) error
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

// NewPermissions ctor
func NewPermissions(psRepo repositories.Permissions) Permissions {
	return &permissions{
		permissionsRepository: psRepo,
	}
}
