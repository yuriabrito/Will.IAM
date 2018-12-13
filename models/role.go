package models

// Role type
type Role struct {
	ID   string `json:"id" pg:"id"`
	Name string `json:"name" pg:"name"`
	// Should change updatedAt when a permission is created for role
	CreatedUpdatedAt
}

// RoleBinding type
type RoleBinding struct {
	ID               string `json:"id" pg:"id"`
	ServiceAccountID string `json:"serviceAccountId" pg:"service_account_id"`
	RoleID           string `json:"roleId" pg:"role_id"`
	CreatedUpdatedAt
}
