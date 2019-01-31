package repositories

// All holds a reference to each possible repository interface
type All struct {
	Permissions     Permissions
	Roles           Roles
	ServiceAccounts ServiceAccounts
	Services        Services
	Tokens          Tokens
	Healthcheck     Healthcheck
	storage         *Storage
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

// WithPGTx clones All and all its contents and injects a PG tx
// in it's storage.PG.DB and in all inner repo storages
func (a *All) WithPGTx() (*All, error) {
	c := &All{
		Permissions:     a.Permissions.Clone(),
		Roles:           a.Roles.Clone(),
		ServiceAccounts: a.ServiceAccounts.Clone(),
		Tokens:          a.Tokens.Clone(),
		storage:         a.storage.Clone(),
	}
	tx, err := a.storage.PG.Begin(a.storage.PG.DB)
	if err != nil {
		return nil, err
	}
	a.storage.PG.DB = tx
	c.Permissions.setStorage(a.storage)
	c.Roles.setStorage(a.storage)
	c.ServiceAccounts.setStorage(a.storage)
	c.Tokens.setStorage(a.storage)
	return c, nil
}
