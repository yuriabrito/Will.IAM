package repositories

import "github.com/ghostec/Will.IAM/models"

// ServiceAccounts repository
type ServiceAccounts interface {
	Get(string) (*models.ServiceAccount, error)
}

type serviceAccounts struct{}

func (s serviceAccounts) Get(id string) (*models.ServiceAccount, error) {
	return nil, nil
}

// NewServiceAccounts serviceAccounts ctor
func NewServiceAccounts() ServiceAccounts {
	return &serviceAccounts{}
}
