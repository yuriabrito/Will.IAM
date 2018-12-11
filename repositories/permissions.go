package repositories

import "github.com/ghostec/Will.IAM/models"

// Permissions repository
type Permissions interface {
	ForServiceAccount(string) ([]models.Permission, error)
}

type permissions struct {
	storage *Storage
}

func (p *permissions) ForServiceAccount(
	id string,
) ([]models.Permission, error) {
	return nil, nil
}

// NewPermissions users ctor
func NewPermissions(s *Storage) Permissions {
	return &permissions{storage: s}
}
