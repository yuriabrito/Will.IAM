package repositories

import (
	"fmt"

	"github.com/ghostec/Will.IAM/models"
)

// Roles repository
type Roles interface {
	GetServiceAccounts(string) ([]models.ServiceAccount, error)
	ForServiceAccountID(string) ([]models.Role, error)
	Create(*models.Role) error
	Update(*models.Role) error
	Bind(models.Role, models.ServiceAccount) error
	WithNamePrefix(string, int) ([]models.Role, error)
	List() ([]models.Role, error)
	Get(string) (*models.Role, error)
}

type roles struct {
	storage *Storage
}

func (rs roles) GetServiceAccounts(
	roleID string,
) ([]models.ServiceAccount, error) {
	sas := []models.ServiceAccount{}
	_, err := rs.storage.PG.DB.Query(
		&sas,
		`SELECT sa.id, sa.name, sa.picture, sa.email FROM service_accounts sa
		JOIN role_bindings rb ON rb.service_account_id = sa.id
		WHERE rb.role_id = ?
		ORDER BY sa.created_at DESC`,
		roleID,
	)
	if err != nil {
		return nil, err
	}
	return sas, nil
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

func (rs roles) Update(r *models.Role) error {
	// tx, err := rs.storage.PG.Begin(rs.storage.PG.DB)
	_, err := rs.storage.PG.DB.Query(
		r, `UPDATE roles SET name = ?name WHERE id = ?id`, r,
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

func (rs roles) WithNamePrefix(
	prefix string, maxResults int,
) ([]models.Role, error) {
	var rsSl []models.Role
	if _, err := rs.storage.PG.DB.Query(
		&rsSl, "SELECT id, name FROM roles WHERE name ILIKE ?",
		fmt.Sprintf("%s%%", prefix),
	); err != nil {
		return nil, err
	}
	return rsSl, nil
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

func (rs roles) Get(id string) (*models.Role, error) {
	r := new(models.Role)
	if _, err := rs.storage.PG.DB.Query(
		r, `SELECT id, name, is_base_role, created_at, updated_at
		FROM roles WHERE id = ?`, id,
	); err != nil {
		return nil, err
	}
	if r.ID == "" {
		return nil, fmt.Errorf("role %s not found", id)
	}
	return r, nil
}

// NewRoles roles ctor
func NewRoles(s *Storage) Roles {
	return &roles{storage: s}
}
