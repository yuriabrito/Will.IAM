package repositories

import (
	"github.com/ghostec/Will.IAM/models"
)

// Roles repository
type Roles interface {
	ForServiceAccountID(string) ([]models.Role, error)
	Create(*models.Role) error
	Bind(models.Role, models.ServiceAccount) error
	List() ([]models.Role, error)
}

type roles struct {
	storage *Storage
}

func (rs roles) ForServiceAccountID(serviceAccountID string) ([]models.Role, error) {
	var roles []models.Role
	_, err := rs.storage.PG.DB.Query(
		&roles,
		`SELECT r.id, r.name FROM roles r
		JOIN role_bindings rb ON rb.role_id = r.id
		WHERE rb.service_account_id = ?`,
		serviceAccountID,
	)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (rs roles) Create(r *models.Role) error {
	_, err := rs.storage.PG.DB.Query(
		r, `INSERT INTO roles (name, is_base_role) VALUES (?name, ?is_base_role)
		RETURNING id`, r,
	)
	return err
}

func (rs roles) Bind(r models.Role, sa models.ServiceAccount) error {
	rb := &models.RoleBinding{
		RoleID:           r.ID,
		ServiceAccountID: sa.ID,
	}
	_, err := rs.storage.PG.DB.Exec(
		`INSERT INTO role_bindings (role_id, service_account_id)
		VALUES (?role_id, ?service_account_id)`, rb,
	)
	return err
}

func (rs roles) List() ([]models.Role, error) {
	var rsSl []models.Role
	if _, err := rs.storage.PG.DB.Query(
		&rsSl, "SELECT id, name FROM roles WHERE is_base_role = false",
	); err != nil {
		return nil, err
	}
	return rsSl, nil
}

// NewRoles roles ctor
func NewRoles(s *Storage) Roles {
	return &roles{storage: s}
}
