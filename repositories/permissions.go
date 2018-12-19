package repositories

import (
	"github.com/ghostec/Will.IAM/models"
	"github.com/go-pg/pg"
)

// Permissions repository
type Permissions interface {
	ForRoles([]models.Role) ([]models.Permission, error)
	Create(models.Permission) error
}

type permissions struct {
	storage *Storage
}

func (ps *permissions) ForRoles(
	roles []models.Role,
) ([]models.Permission, error) {
	rolesIds := make([]string, len(roles))
	for i := range roles {
		rolesIds[i] = roles[i].ID
	}

	var permissions []models.Permission
	if _, err := ps.storage.PG.DB.Query(
		&permissions, `SELECT role_id, service, ownership_level,
action, resource_hierarchy FROM permissions
	WHERE role_id = ANY (?)`, pg.Array(rolesIds),
	); err != nil {
		return nil, err
	}
	return permissions, nil
}

func (ps *permissions) Create(p models.Permission) error {
	_, err := ps.storage.PG.DB.Exec(
		`INSERT INTO permissions (role_id, service, ownership_level, action,
		resource_hierarchy) VALUES (?, ?, ?, ?, ?)`, p.RoleID, p.Service,
		p.OwnershipLevel, p.Action, p.ResourceHierarchy,
	)
	return err
}

// NewPermissions users ctor
func NewPermissions(s *Storage) Permissions {
	return &permissions{storage: s}
}
