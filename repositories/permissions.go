package repositories

// Permissions repository
type Permissions interface{}

type permissions struct{}

// NewPermissions users ctor
func NewPermissions() Permissions {
	return &permissions{}
}
