package repositories

import (
	"fmt"

	"github.com/ghostec/Will.IAM/errors"
	"github.com/ghostec/Will.IAM/models"
)

// ServiceAccounts repository
type ServiceAccounts interface {
	Get(string) (*models.ServiceAccount, error)
	List(*ListOptions) ([]models.ServiceAccount, error)
	Count() (int64, error)
	Search(string) ([]models.ServiceAccount, error)
	ForEmail(string) (*models.ServiceAccount, error)
	ForKeyPair(string, string) (*models.ServiceAccount, error)
	Create(*models.ServiceAccount) error
	Update(*models.ServiceAccount) error
	Clone() ServiceAccounts
	setStorage(*Storage)
}

type serviceAccounts struct {
	*withStorage
}

func (sas *serviceAccounts) Clone() ServiceAccounts {
	return NewServiceAccounts(sas.storage.Clone())
}

func (sas serviceAccounts) Get(id string) (*models.ServiceAccount, error) {
	sa := new(models.ServiceAccount)
	if _, err := sas.storage.PG.DB.Query(
		sa,
		`SELECT id, name, key_id, key_secret, email, base_role_id, picture
		FROM service_accounts
		WHERE id = ?`,
		id,
	); err != nil {
		return nil, err
	}
	if sa.ID == "" {
		return nil, errors.NewEntityNotFoundError(models.ServiceAccount{}, id)
	}
	if sa.KeyID != "" {
		sa.AuthenticationType = models.AuthenticationTypes.KeyPair
	}
	return sa, nil
}

func (sas serviceAccounts) List(
	lo *ListOptions,
) ([]models.ServiceAccount, error) {
	var saSl []models.ServiceAccount
	if _, err := sas.storage.PG.DB.Query(
		&saSl,
		`SELECT id, name, email, picture FROM service_accounts
		ORDER BY name ASC LIMIT ? OFFSET ?`, lo.PageSize, lo.Page*lo.PageSize,
	); err != nil {
		return nil, err
	}
	for i := range saSl {
		authType := models.AuthenticationTypes.OAuth2
		if saSl[i].KeyID != "" {
			authType = models.AuthenticationTypes.KeyPair
		}
		saSl[i].AuthenticationType = authType
	}
	return saSl, nil
}

func (sas serviceAccounts) Count() (int64, error) {
	var count int64
	if _, err := sas.storage.PG.DB.Query(
		&count,
		`SELECT count(*) FROM service_accounts`,
	); err != nil {
		return 0, err
	}
	return count, nil
}

func (sas serviceAccounts) Search(
	term string,
) ([]models.ServiceAccount, error) {
	saSl := []models.ServiceAccount{}
	if _, err := sas.storage.PG.DB.Query(
		&saSl,
		`SELECT id, name, email, picture FROM service_accounts
		WHERE name ILIKE ?0 OR email ILIKE ?0
		ORDER BY name ASC`,
		fmt.Sprintf("%%%s%%", term),
	); err != nil {
		return nil, err
	}
	for i := range saSl {
		authType := models.AuthenticationTypes.OAuth2
		if saSl[i].KeyID != "" {
			authType = models.AuthenticationTypes.KeyPair
		}
		saSl[i].AuthenticationType = authType
	}
	return saSl, nil
}

// ForEmail retrieves Service Account corresponding
func (sas serviceAccounts) ForEmail(
	email string,
) (*models.ServiceAccount, error) {
	sa := new(models.ServiceAccount)
	if _, err := sas.storage.PG.DB.Query(
		sa, `SELECT id, name, key_id, key_secret, email, base_role_id, picture
		FROM service_accounts WHERE email = ? LIMIT 1`, email,
	); err != nil {
		return nil, err
	}
	if sa.ID == "" {
		return nil, errors.NewEntityNotFoundError(models.ServiceAccount{}, email)
	}
	return sa, nil
}

// ForKeyPair retrieves Service Account corresponding
func (sas serviceAccounts) ForKeyPair(
	keyID, keySecret string,
) (*models.ServiceAccount, error) {
	sa := []*models.ServiceAccount{}
	if _, err := sas.storage.PG.DB.Query(
		&sa, `SELECT id, name, key_id, key_secret, email, base_role_id
		FROM service_accounts WHERE key_id = ? AND key_secret = ?`,
		keyID, keySecret,
	); err != nil {
		return nil, err
	}
	if len(sa) == 0 {
		return nil, errors.NewEntityNotFoundError(models.ServiceAccount{}, keyID)
	}
	return sa[0], nil
}

func (sas serviceAccounts) Create(sa *models.ServiceAccount) error {
	_, err := sas.storage.PG.DB.Query(
		sa, `INSERT INTO service_accounts (id, name, email, key_id, key_secret,
		base_role_id) VALUES (?id, ?name, ?email, ?key_id, ?key_secret,
		?base_role_id) RETURNING id`, sa,
	)
	return err
}

func (sas serviceAccounts) Update(sa *models.ServiceAccount) error {
	_, err := sas.storage.PG.DB.Exec(
		`UPDATE service_accounts SET name = ?name, email = ?email,
		key_id = ?key_id, key_secret = ?key_secret, base_role_id = ?base_role_id,
		picture = ?picture, updated_at = now() WHERE id = ?id`, sa,
	)
	return err
}

// NewServiceAccounts serviceAccounts ctor
func NewServiceAccounts(s *Storage) ServiceAccounts {
	return &serviceAccounts{&withStorage{storage: s}}
}
