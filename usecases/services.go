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
	servicesRepository     repositories.Services
	serviceAccountsUseCase ServiceAccounts
}

// Create a new service with unique name and permission name
// Also creates an associate Service Account with full access
// and attributes full access to creator
func (ss services) Create(
	service *models.Service, creatorServiceAccountID string,
) error {
	// TODO: use tx
	sa := models.BuildKeyPairServiceAccount(service.Name)
	if err := ss.serviceAccountsUseCase.Create(sa); err != nil {
		return err
	}
	service.ServiceAccountID = sa.ID
	if err := ss.servicesRepository.Create(service); err != nil {
		return err
	}
	fullAccessPermission := models.Permission{
		Service:           service.PermissionName,
		OwnershipLevel:    models.OwnershipLevels.Owner,
		Action:            models.Action("*"),
		ResourceHierarchy: models.ResourceHierarchy("*"),
	}
	ss.serviceAccountsUseCase.CreatePermission(
		sa.ID, fullAccessPermission,
	)
	ss.serviceAccountsUseCase.CreatePermission(
		creatorServiceAccountID, fullAccessPermission,
	)
	return nil
}

// NewServices services' ctor
func NewServices(
	servicesRepository repositories.Services,
	serviceAccountsUseCase ServiceAccounts,
) Services {
	return &services{
		servicesRepository:     servicesRepository,
		serviceAccountsUseCase: serviceAccountsUseCase,
	}
}
