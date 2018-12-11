package models

// RolePermissions type
type RolePermissions struct {
	ID           string `json:"id" pg:"id"`
	RoleID       string `json:"roleId" pg:"role_id"`
	PermissionID string `json:"permissionId" pg:"permission_id"`
	CreatedUpdatedAt
}
