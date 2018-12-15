package repositories

import (
	"fmt"
	"strings"

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

	type pgPermission struct {
		RoleID            string   `pg:"role_id"`
		Service           string   `pg:"service"`
		OwnershipLevel    string   `pg:"ownership_level"`
		Action            string   `pg:"action"`
		ResourceHierarchy []string `pg:"resource_hierarchy,array"`
	}

	var pgPss []pgPermission
	_, err := ps.storage.PG.DB.Query(
		&pgPss, `SELECT role_id, service, ownership_level,
action, resource_hierarchy FROM permissions
	WHERE role_id = ANY (?)`, pg.Array(rolesIds),
	)
	if err != nil {
		return nil, err
	}
	pss := make([]models.Permission, len(pgPss))
	for i := range pgPss {
		if pss[i], err = models.BuildPermission(fmt.Sprintf(
			"%s::%s::%s::%s", pgPss[i].Service, pgPss[i].OwnershipLevel,
			pgPss[i].Action, strings.Join(pgPss[i].ResourceHierarchy, "::"),
		)); err != nil {
			return nil, err
		}
		pss[i].RoleID = pgPss[i].RoleID
	}
	return pss, nil
}

func (ps *permissions) Create(p models.Permission) error {
	_, err := ps.storage.PG.DB.Exec(
		`INSERT INTO permissions (role_id, service, ownership_level, action,
		resource_hierarchy) VALUES (?, ?, ?, ?, ?)`, p.RoleID, p.Service,
		p.OwnershipLevel, p.Action, pg.Array(p.ResourceHierarchy.Hierarchy),
	)
	return err
}

// NewPermissions users ctor
func NewPermissions(s *Storage) Permissions {
	return &permissions{storage: s}
}
