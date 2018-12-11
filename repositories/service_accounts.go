package repositories

import "github.com/ghostec/Will.IAM/models"

// ServiceAccounts repository
type ServiceAccounts interface {
	Get(string) (*models.ServiceAccount, error)
	Create(*models.ServiceAccount) error
}

type serviceAccounts struct {
	storage *Storage
}

func (s serviceAccounts) Get(id string) (*models.ServiceAccount, error) {
	return nil, nil
}

func (s serviceAccounts) Create(sa *models.ServiceAccount) error {
	_, err := s.storage.PG.DB.Query(
		sa, "INSERT INTO service_accounts (email) VALUES (?email) RETURNING id", sa,
	)
	return err
}

// NewServiceAccounts serviceAccounts ctor
func NewServiceAccounts(s *Storage) ServiceAccounts {
	return &serviceAccounts{storage: s}
}
