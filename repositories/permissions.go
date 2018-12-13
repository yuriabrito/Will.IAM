package repositories

import "github.com/ghostec/Will.IAM/models"

// Permissions repository
type Permissions interface {
	ForRoles([]models.Role) ([]models.Permission, error)
}

type permissions struct {
	storage *Storage
}

func (p *permissions) ForRoles(
	roles []models.Role,
) ([]models.Permission, error) {
	return nil, nil
}

// NewPermissions users ctor
func NewPermissions(s *Storage) Permissions {
	return &permissions{storage: s}
}
