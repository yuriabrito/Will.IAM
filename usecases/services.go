package usecases

import (
	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/repositories"
)

// Services contract
type Services interface {
	All() ([]models.Service, error)
	Create(*models.Service, string) error
}

type services struct {
	repo *repositories.All
}

// Create a new service with unique name and permission name
// Also creates an associate Service Account with full access
// and attributes full access to creator
func (ss services) Create(
	service *models.Service, creatorServiceAccountID string,
) error {
	creatorSA, err := ss.repo.ServiceAccounts.Get(creatorServiceAccountID)
	if err != nil {
		return err
	}
	return ss.repo.WithPGTx(func(repo *repositories.All) error {
		sa := models.BuildKeyPairServiceAccount(service.Name)
		if err := createServiceAccount(sa, repo); err != nil {
			return err
		}
		service.ServiceAccountID = sa.ID
		if err := ss.repo.Services.Create(service); err != nil {
			return err
		}
		buildFullAccessPermissionForRoleID := func(
			roleID string,
		) *models.Permission {
			return &models.Permission{
				Service:           service.PermissionName,
				OwnershipLevel:    models.OwnershipLevels.Owner,
				Action:            models.Action("*"),
				ResourceHierarchy: models.ResourceHierarchy("*"),
				RoleID:            roleID,
			}
		}
		// TODO: check if base_role_id is set in sa
		repo.Permissions.Create(
			buildFullAccessPermissionForRoleID(sa.BaseRoleID),
		)
		repo.Permissions.Create(
			buildFullAccessPermissionForRoleID(creatorSA.BaseRoleID),
		)
		return nil
	})
}

func (ss services) All() ([]models.Service, error) {
	return ss.repo.Services.All()
}

// NewServices services' ctor
func NewServices(repo *repositories.All) Services {
	return &services{repo: repo}
}
