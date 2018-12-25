package repositories

import (
	"fmt"

	"github.com/ghostec/Will.IAM/models"
	"github.com/go-pg/pg"
)

// Permissions repository
type Permissions interface {
	Get(string) (*models.Permission, error)
	ForRoles([]models.Role) ([]models.Permission, error)
	Create(*models.Permission) error
	Delete(string) error
}

type permissions struct {
	storage *Storage
}

// Get retrieves a permission by id
func (ps *permissions) Get(id string) (*models.Permission, error) {
	p := new(models.Permission)
	if info, err := ps.storage.PG.DB.Query(
		p, `SELECT id, role_id, service, ownership_level,
action, resource_hierarchy FROM permissions
	WHERE id = ?`, id,
	); err != nil {
		return nil, err
	} else if info.RowsReturned() == 0 {
		return nil, fmt.Errorf("permission %s not found", id)
	}
	return p, nil
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
		&permissions, `SELECT id, role_id, service, ownership_level,
action, resource_hierarchy FROM permissions
	WHERE role_id = ANY (?)`, pg.Array(rolesIds),
	); err != nil {
		return nil, err
	}
	return permissions, nil
}

func (ps *permissions) Create(p *models.Permission) error {
	_, err := ps.storage.PG.DB.Exec(
		`INSERT INTO permissions (role_id, service, ownership_level, action,
		resource_hierarchy) VALUES (?, ?, ?, ?, ?) RETURNING id`, p.RoleID,
		p.Service, p.OwnershipLevel, p.Action, p.ResourceHierarchy,
	)
	return err
}

func (ps *permissions) Delete(id string) error {
	_, err := ps.storage.PG.DB.Exec(
		`DELETE FROM permissions WHERE id = ?`, id,
	)
	return err
}

// NewPermissions users ctor
func NewPermissions(s *Storage) Permissions {
	return &permissions{storage: s}
}
