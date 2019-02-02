package repositories

import (
	"context"

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

// WithCtx clones All and all its contents and injects a context
// in all storages
func (a *All) WithCtx(ctx context.Context) *All {
	s := a.storage.Clone()
	s.PG.DB = s.PG.DB.WithContext(ctx)
	s.Redis.Client = s.Redis.Client.WithContext(ctx)
	return a.cloneWithStorage(s)
}

// WithPGTx clones All and all its contents and injects a PG tx
// in it's storage.PG.DB and in all inner repo storages
func (a *All) WithPGTx(ctx context.Context, fn func(repo *All) error) error {
	s := a.storage.Clone()
	tx, err := a.storage.PG.Begin(a.storage.PG.DB.WithContext(ctx))
	if err != nil {
		return err
	}
	s.PG.DB = tx
	s.Redis.Client = s.Redis.Client.WithContext(ctx)
	c := a.cloneWithStorage(s)

	defer pg.Rollback(c.storage.PG.DB)
	err = fn(c)
	if err == nil {
		return pg.Commit(c.storage.PG.DB)
	}
	return nil
}

func (a *All) cloneWithStorage(s *Storage) *All {
	c := &All{
		Permissions:     a.Permissions.Clone(),
		Roles:           a.Roles.Clone(),
		ServiceAccounts: a.ServiceAccounts.Clone(),
		Services:        a.Services.Clone(),
		Tokens:          a.Tokens.Clone(),
		storage:         s,
	}
	c.Permissions.setStorage(s)
	c.Roles.setStorage(s)
	c.ServiceAccounts.setStorage(s)
	c.Services.setStorage(s)
	c.Tokens.setStorage(s)
	return c
}
