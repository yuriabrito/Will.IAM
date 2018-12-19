package usecases

import (
	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/repositories"
)

// Services contract
type Services interface {
	Create(*models.Service, string) error
}

type services struct {
	servicesRepository repositories.Services
}

// Create a new service with unique name and permission name
// Also creates an associate Service Account with full access
// and attributes full access to creator
func (ss services) Create(
	service *models.Service, creatorServiceAccountID string,
) error {
	// TODO: use tx
	return ss.servicesRepository.Create(service)
}

// NewServices services' ctor
func NewServices(servicesRepository repositories.Services) Services {
	return &services{servicesRepository: servicesRepository}
}
