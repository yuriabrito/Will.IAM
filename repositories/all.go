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
func (a All) WithPGTx() (*All, error) {
	c := &All{
		Permissions:     a.Permissions.Clone(),
		Roles:           a.Roles.Clone(),
		ServiceAccounts: a.ServiceAccounts.Clone(),
		Services:        a.Services.Clone(),
		Tokens:          a.Tokens.Clone(),
		storage:         a.storage.Clone(),
	}
	tx, err := c.storage.PG.Begin(c.storage.PG.DB)
	if err != nil {
		return nil, err
	}
	c.storage.PG.DB = tx
	c.Permissions.setStorage(c.storage)
	c.Roles.setStorage(c.storage)
	c.ServiceAccounts.setStorage(c.storage)
	c.Services.setStorage(c.storage)
	c.Tokens.setStorage(c.storage)
	return c, nil
}

// PGTxCommit commits the tx in a.storage.PG.DB
func (a *All) PGTxCommit() error {
	return pg.Commit(a.storage.PG.DB)
}

// PGTxRollback does a rollback to the tx in a.storage.PG.DB
func (a *All) PGTxRollback() error {
	return pg.Rollback(a.storage.PG.DB)
}
