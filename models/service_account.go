package models

import "github.com/gofrs/uuid"

// ServiceAccount type
type ServiceAccount struct {
	ID         string `json:"id" pg:"id"`
	Name       string `json:"name" pg:"name"`
	KeyID      string `json:"keyId" pg:"key_id"`
	KeySecret  string `json:"keySecret" pg:"key_secret"`
	Email      string `json:"email" pg:"email"`
	BaseRoleID string `json:"baseRoleId" pg:"base_role_id"`
	CreatedUpdatedAt
}

// BuildKeyPairServiceAccount generates random KeyID and KeySecret
func BuildKeyPairServiceAccount(name string) *ServiceAccount {
	return &ServiceAccount{
		Name:      name,
		KeyID:     uuid.Must(uuid.NewV4()).String(),
		KeySecret: uuid.Must(uuid.NewV4()).String(),
	}
}

// BuildOAuth2ServiceAccount generates a ServiceAccount with Name and Email
func BuildOAuth2ServiceAccount(name, email string) *ServiceAccount {
	return &ServiceAccount{
		Name:  name,
		Email: email,
	}
}
