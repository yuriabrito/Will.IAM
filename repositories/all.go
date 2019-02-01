package repositories

import (
	"github.com/topfreegames/extensions/pg"
)

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
		Services:        NewServices(s),
		Tokens:          NewTokens(s),
		Healthcheck:     NewHealthcheck(s),
		storage:         s,
	}
}

// WithPGTx clones All and all its contents and injects a PG tx
// in it's storage.PG.DB and in all inner repo storages
func (a *All) WithPGTx(fn func(repo *All) error) error {
	c := &All{
		Permissions:     a.Permissions.Clone(),
		Roles:           a.Roles.Clone(),
		ServiceAccounts: a.ServiceAccounts.Clone(),
		Services:        a.Services.Clone(),
		Tokens:          a.Tokens.Clone(),
		storage:         a.storage.Clone(),
	}
	tx, err := a.storage.PG.Begin(a.storage.PG.DB)
	if err != nil {
		return err
	}
	c.storage.PG.DB = tx
	c.Permissions.setStorage(c.storage)
	c.Roles.setStorage(c.storage)
	c.ServiceAccounts.setStorage(c.storage)
	c.Services.setStorage(c.storage)
	c.Tokens.setStorage(c.storage)

	defer pg.Rollback(c.storage.PG.DB)
	err = fn(c)
	if err == nil {
		return pg.Commit(c.storage.PG.DB)
	}
	return nil
}
