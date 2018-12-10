package repositories

import "github.com/ghostec/Will.IAM/models"

// Permissions repository
type Permissions interface {
	ForServiceAccount(string) ([]models.Permission, error)
}

type permissions struct{}

func (p *permissions) ForServiceAccount(
	id string,
) ([]models.Permission, error) {
	return nil, nil
}

// NewPermissions users ctor
func NewPermissions() Permissions {
	return &permissions{}
}
