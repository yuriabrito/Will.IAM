package models

// Service type
type Service struct {
	ID                      string `json:"id" pg:"id"`
	Name                    string `json:"name" pg:"name"`
	PermissionName          string `json:"permissionName" pg:"permission_name"`
	CreatorServiceAccountID string `json:"creatorServiceAccountID" pg:"creator_service_account_id"`
	CreatedUpdatedAt
}
