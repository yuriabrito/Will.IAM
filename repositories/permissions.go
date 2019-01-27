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
	CreateRequest(string, *models.PermissionRequest) error
	GetPermissionRequests(string) ([]models.PermissionRequest, error)
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

// GetPermissionRequets retrieve permission requests for a service account
func (ps *permissions) GetPermissionRequests(
	saID string,
) ([]models.PermissionRequest, error) {
	prs := []models.PermissionRequest{}
	_, err := ps.storage.PG.DB.Query(
		&prs, `SELECT id, service, action, resource_hierarchy, message, state,
		created_at, updated_at FROM permissions_requests
		WHERE service_account_id = ?`, saID,
	)
	if err != nil {
		return nil, err
	}
	return prs, nil
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
		resource_hierarchy) VALUES (?, ?, ?, ?, ?)
		ON CONFLICT DO NOTHING RETURNING id`, p.RoleID,
		p.Service, p.OwnershipLevel, p.Action, p.ResourceHierarchy,
	)
	return err
}
func (ps *permissions) CreateRequest(
	saID string, r *models.PermissionRequest,
) error {
	_, err := ps.storage.PG.DB.Query(r,
		`INSERT INTO permissions_requests (service, action, resource_hierarchy,
		message, state, service_account_id) VALUES (?, ?, ?, ?, ?, ?) RETURNING id`,
		r.Service, r.Action, r.ResourceHierarchy, r.Message, r.State, saID)
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
