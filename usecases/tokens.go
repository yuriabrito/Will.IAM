package usecases

import "github.com/ghostec/Will.IAM/models"

// Tokens contract
type Tokens interface {
	IsValid() (bool, error)
	Refresh() error
	Authenticate(models.ServiceAccount) (bool, error)
}
