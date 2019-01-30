package repositories

// All holds a reference to each possible repository interface
type All struct {
	Permissions     Permissions
	Roles           Roles
	ServiceAccounts ServiceAccounts
	Services        Services
	Tokens          Tokens
	Healthcheck     Healthcheck
}

// New All ctor
func New(s *Storage) *All {
	return &All{
		Permissions:     NewPermissions(s),
		Roles:           NewRoles(s),
		ServiceAccounts: NewServiceAccounts(s),
		Tokens:          NewTokens(s),
		Healthcheck:     NewHealthcheck(s),
	}
}
